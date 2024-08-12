package ai

import (
	discussion_model "code.gitea.io/gitea/models/discussion"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"fmt"
	"github.com/spf13/cast"
	"sync"
)

type DiscussionAiService interface {
	GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm, aiRequester AiSampleCodeRequester, adapter DiscussionDbAdapter) ([]*AiSampleCodeResponse, error)
	CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm, adapter DiscussionDbAdapter) (*discussion_model.DiscussionAiSampleCode, error)
	DeleteAiSampleCode(ctx *context.Context, id int64, adapter DiscussionDbAdapter) error
}

type DiscussionAiServiceImpl struct{}

var _ DiscussionAiService = &DiscussionAiServiceImpl{}

func (is *DiscussionAiServiceImpl) CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm, adapter DiscussionDbAdapter) (*discussion_model.DiscussionAiSampleCode, error) {
	// TODOC Discussion Comment 무결성 검사

	aiSampleCode, err := adapter.InsertDiscussionAiSampleCode(ctx, &discussion_model.CreateDiscussionAiCommentOpt{
		TargetCommentId: cast.ToInt64(form.TargetCommentId),
		GenearaterId:    ctx.Doer.ID,
		Content:         &form.SampleCodeContent,
	})

	if err != nil {
		return nil, err
	}

	return aiSampleCode, nil

}

// TODOC 재시도 횟수를 유저 정보로부터 가져와서 제한하기.
// TODOC 잘못된 형식의 json이 돌아올 때 예외 반환하기(json 형식 표시하도록)
func (is *DiscussionAiServiceImpl) GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm, aiRequester AiSampleCodeRequester, adapter DiscussionDbAdapter) ([]*AiSampleCodeResponse, error) {
	wg := new(sync.WaitGroup)

	wg.Add(AI_SAMPLE_CODE_UNIT)

	resultQueue := make(chan *AiSampleCodeResponse, AI_SAMPLE_CODE_UNIT)

	for i := 0; i < AI_SAMPLE_CODE_UNIT; i++ {
		go func(form *api.GenerateAiSampleCodesForm) {
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

// TODOC ai 샘플코드 저장

// TODOC ai 샘플코드 삭제

// TODOC a

func (aiService *DiscussionAiServiceImpl) DeleteAiSampleCode(ctx *context.Context, id int64, adapter DiscussionDbAdapter) error {

	return adapter.DeleteDiscussionAiSampleCodeByID(ctx, id)
}
