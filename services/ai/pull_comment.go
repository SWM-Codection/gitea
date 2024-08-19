package ai

import (
	"fmt"
	"strconv"
	"sync"

	issues_model "code.gitea.io/gitea/models/issues"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

type PullCommentService interface {
	CreateAiPullComment(ctx *context.Context, form *api.CreateAiPullCommentForm) error
	DeleteAiPullComment(ctx *context.Context, id int64) error
}

var _ PullCommentService = &PullCommentServiceImpl{}

// TODOC 잘못된 형식의 json이 돌아올 때 예외 반환하기(json 형식 표시하도록)
func (is *PullCommentServiceImpl) CreateAiPullComment(ctx *context.Context, form *api.CreateAiPullCommentForm) error {
	branch := form.Branch
	wg := new(sync.WaitGroup)

	requestCnt := len(*form.FileContents)
	wg.Add(requestCnt)

	resultQueue := make(chan *AiPullCommentResponse, requestCnt)

	for _, fileContent := range *form.FileContents {
		go func(fileContent *api.PathContentMap) {
			defer wg.Done()
			result, err := AiPullCommentRequester.RequestReviewToAI(ctx, &AiPullCommentRequest{
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

	return saveResults(ctx, resultQueue, pullID)
}

func (is *PullCommentServiceImpl) DeleteAiPullComment(ctx *context.Context, id int64) error {

	return AiPullCommentDbAdapter.DeleteAiPullCommentByID(ctx, id)
}

func saveResults(ctx *context.Context, reviewResults chan *AiPullCommentResponse, pullID int64) error {
	pull, err := AiPullCommentDbAdapter.GetIssueByID(ctx, pullID)
	if err != nil {
		return fmt.Errorf("pr not found by id")
	}

	for result := range reviewResults {
		_, err := AiPullCommentDbAdapter.CreateAiPullComment(ctx, &issues_model.CreateAiPullCommentOption{
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

type PullCommentDbAdapter interface {
	GetIssueByID(ctx *context.Context, id int64) (*issues_model.Issue, error)
	CreateAiPullComment(ctx *context.Context, opts *issues_model.CreateAiPullCommentOption) (*issues_model.AiPullComment, error)
	DeleteAiPullCommentByID(ctx *context.Context, id int64) error
}

type PullCommentDbAdapterImpl struct{}

func (is *PullCommentDbAdapterImpl) GetIssueByID(ctx *context.Context, id int64) (*issues_model.Issue, error) {
	return issues_model.GetIssueByID(ctx, id)
}

func (is *PullCommentDbAdapterImpl) CreateAiPullComment(ctx *context.Context, opts *issues_model.CreateAiPullCommentOption) (*issues_model.AiPullComment, error) {
	return issues_model.CreateAiPullComment(ctx, opts)
}

func (is *PullCommentDbAdapterImpl) DeleteAiPullCommentByID(ctx *context.Context, id int64) error {

	return issues_model.DeleteAiPullCommentByID(ctx, id)

}
