package ai

import (
	"fmt"
	"sync"

	ai_client "code.gitea.io/gitea/client/ai"
	discussion_model "code.gitea.io/gitea/models/discussion"
	"code.gitea.io/gitea/modules/json"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
)

const AI_SAMPLE_CODE_UNIT = 3

type DiscussionAiService interface {
	GenerateAiSampleCodes(ctx *context.Context, form *api.CreateSampleAiCommentForm, aiRequester AiSampleCodeRequester, adapter DiscussionDbAdapter) ([]*AiSampleCodeResponse, error)
	DeleteAiPullComment(ctx *context.Context, id int64, adapter DiscussionDbAdapter) error
}

type AiSampleCodeRequest struct {
	CommentContent string `json:"comment"`
	CodeContent    string `json:"code"`
}

type AiSampleCodeResponse struct {
	SampleCode string `json:"sample_code"`
}

type DiscussionAiServiceImpl struct{}

type AiSampleCodeRequesterImpl struct{}

type AiSampleCodeRequester interface {
	RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error)
}

var _ DiscussionAiService = &DiscussionAiServiceImpl{}
var _ AiSampleCodeRequester = &AiSampleCodeRequesterImpl{}

type AiSampleCodeResult struct {
	SampleCode string `json:"sample_code"`
}

func (aiRequest *AiSampleCodeRequesterImpl) RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error) {

	response, err := ai_client.Request().SetBody(request).Post("/api/sample")

	if err != nil {
		fmt.Errorf("%s", err.Error())

		return nil, err
	}

	result := &AiSampleCodeResponse{}
	err = json.Unmarshal(response.Body(), result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func saveAiSampleCode(ctx *context.Context, reviewResults chan *AiSampleCodeResponse, pullID int64, adapter DiscussionDbAdapter) error {

	// TODOC AI 샘플 코드를 받아서 저장하기

	return nil
}

type DiscussionDbAdapter interface {
	// GetDiscussionCommentByID(ctx *context.Context, id int64) (*issues_model.Issue, error)
	CreateDiscussionAiComment(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiComment, error)
	DeleteDiscussionAiCommentByID(ctx *context.Context, id int64) error
}

type DiscussionDbAdapterImpl struct{}

func (is *DiscussionDbAdapterImpl) CreateDiscussionAiComment(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiComment, error) {
	return discussion_model.CreateDiscussionAiComment(ctx, opts)
}

func (is *DiscussionDbAdapterImpl) DeleteDiscussionAiCommentByID(ctx *context.Context, id int64) error {

	return discussion_model.DeleteDiscussionAiCommentByID(ctx, id)

}

// TODOC 재시도 횟수를 유저 정보로부터 가져와서 제한하기.
// TODOC 잘못된 형식의 json이 돌아올 때 예외 반환하기(json 형식 표시하도록)
func (aiController *DiscussionAiServiceImpl) GenerateAiSampleCodes(ctx *context.Context, form *api.CreateSampleAiCommentForm, aiRequester AiSampleCodeRequester, adapter DiscussionDbAdapter) ([]*AiSampleCodeResponse, error) {
	var wg *sync.WaitGroup = new(sync.WaitGroup)

	wg.Add(AI_SAMPLE_CODE_UNIT)

	resultQueue := make(chan *AiSampleCodeResponse, AI_SAMPLE_CODE_UNIT)

	for i := 0; i < AI_SAMPLE_CODE_UNIT; i++ {
		go func(form *api.CreateSampleAiCommentForm) {
			defer wg.Done()
			result, err := aiRequester.RequestReviewToAI(ctx, &AiSampleCodeRequest{
				CodeContent:    form.CodeContent,
				CommentContent: form.CommentContent,
			})

			if err != nil {
				// TODOC 실패 시 재시도 로직 추가
				fmt.Errorf("request sample to ai server fail")
				resultQueue <- nil
				return
			}
			resultQueue <- result
		}(form)
	}

	wg.Wait()
	close(resultQueue)
	sampleCodes := make([]*AiSampleCodeResponse, 0, AI_SAMPLE_CODE_UNIT)
	for result := range resultQueue {
		sampleCodes = append(sampleCodes, result)

	}

	return sampleCodes, nil
}

func (aiService *DiscussionAiServiceImpl) DeleteAiPullComment(ctx *context.Context, id int64, adapter DiscussionDbAdapter) error {

	return adapter.DeleteDiscussionAiCommentByID(ctx, id)
}
