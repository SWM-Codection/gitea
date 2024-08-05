package ai

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"code.gitea.io/gitea/client"
	issues_model "code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

type AiService interface {
	CreateAiPullComment(ctx *context.Context, form *api.CreateAiPullCommentForm, aiRequester AiRequester, adapter DbAdapter) error
	DeleteAiPullComment(ctx *context.Context, id int64, adapter DbAdapter) error
}

type AiReviewRequest struct {
	Branch   string `json:"branch"`
	TreePath string `json:"file_path"`
	Content  string `json:"code"`
}

type AiReviewResponse struct {
	Branch   string `json:"branch"`
	TreePath string `json:"file_path"`
	Content  string `json:"code"`
}

type AiServiceImpl struct{}

type AiRequesterImpl struct{}

type AiRequester interface {
	RequestReviewToAI(ctx *context.Context, request *AiReviewRequest) (*AiReviewResponse, error)
}

var _ AiService = &AiServiceImpl{}
var _ AiRequester = &AiRequesterImpl{}

var apiURL = setting.AiServer.Url

type AiReviewCommentResult struct {
	PullID   string
	RepoID   string
	TreePath string
	Content  string
}

func (aiRequest *AiRequesterImpl) RequestReviewToAI(ctx *context.Context, request *AiReviewRequest) (*AiReviewResponse, error) {

	response, err := client.Request().SetBody(request).Post(fmt.Sprintf(apiURL + "/api/sample"))

	if err != nil {
		return nil, err
	}
	result := &AiReviewResponse{}

	json.Unmarshal(response.Body(), result)

	return result, nil
}

// TODOC 잘못된 형식의 json이 돌아올 때 예외 반환하기(json 형식 표시하도록)
func (aiController *AiServiceImpl) CreateAiPullComment(ctx *context.Context, form *api.CreateAiPullCommentForm, aiRequester AiRequester, adapter DbAdapter) error {
	branch := form.Branch
	var wg *sync.WaitGroup = new(sync.WaitGroup)

	requestCnt := len(*form.FileContents)
	wg.Add(requestCnt)

	resultQueue := make(chan *AiReviewResponse, requestCnt)

	for _, fileContent := range *form.FileContents {
		go func(fileContent *api.PathContentMap) {
			defer wg.Done()
			result, err := aiRequester.RequestReviewToAI(ctx, &AiReviewRequest{
				Branch:   branch,
				TreePath: fileContent.TreePath,
				Content:  fileContent.Content,
			})
			if err != nil {
				fmt.Errorf("request to ai server fail %s", result.TreePath)
				resultQueue <- nil
				return
			}
			resultQueue <- result
		}(&fileContent)
	}

	wg.Wait()
	close(resultQueue)

	pullID, err := strconv.ParseInt(form.PullID, 10, 64)
	if err != nil {
		return fmt.Errorf("pullID is invalid")
	}

	return saveResults(ctx, resultQueue, pullID, adapter)
}

func (aiService *AiServiceImpl) DeleteAiPullComment(ctx *context.Context, id int64, adapter DbAdapter) error {

	return adapter.DeleteAiPullCommentByID(ctx, id)
}

func saveResults(ctx *context.Context, reviewResults chan *AiReviewResponse, pullID int64, adapter DbAdapter) error {
	pull, err := adapter.GetIssueByID(ctx, pullID)
	if err != nil {
		return fmt.Errorf("pr not found by id")
	}

	for result := range reviewResults {
		_, err := adapter.CreateAiPullComment(ctx, &issues_model.CreateAiPullCommentOption{
			Doer:     ctx.Doer,
			Repo:     pull.Repo,
			Pull:     pull,
			TreePath: result.TreePath,
			Content:  result.Content,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type DbAdapter interface {
	GetIssueByID(ctx *context.Context, id int64) (*issues_model.Issue, error)
	CreateAiPullComment(ctx *context.Context, opts *issues_model.CreateAiPullCommentOption) (*issues_model.AiPullComment, error)
	DeleteAiPullCommentByID(ctx *context.Context, id int64) error
}

type DbAdapterImpl struct{}

func (is *DbAdapterImpl) GetIssueByID(ctx *context.Context, id int64) (*issues_model.Issue, error) {
	return issues_model.GetIssueByID(ctx, id)
}

func (is *DbAdapterImpl) CreateAiPullComment(ctx *context.Context, opts *issues_model.CreateAiPullCommentOption) (*issues_model.AiPullComment, error) {
	return issues_model.CreateAiPullComment(ctx, opts)
}

func (is *DbAdapterImpl) DeleteAiPullCommentByID(ctx *context.Context, id int64) error {

	return issues_model.DeleteAiPullCommentByID(ctx, id)

}
