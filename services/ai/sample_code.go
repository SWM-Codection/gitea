package ai

import (
	discussion_model "code.gitea.io/gitea/models/discussion"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"fmt"
	"github.com/spf13/cast"
	"sync"
)

type SampleCodeService interface {
	GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error)
	GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*AiSampleCodeResponse, error)
	CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm) (*discussion_model.AiSampleCode, error)
	DeleteAiSampleCode(ctx *context.Context, id int64) error
}

var DEFAULT_CAPACITY int64 = 10

type SampleCodeServiceImpl struct{}

var _ SampleCodeService = &SampleCodeServiceImpl{}

func (is *SampleCodeServiceImpl) GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error) {

	response, err := AiSampleCodeDbAdapter.GetAiSampleCodesByCommentID(ctx, commentID)

	if err != nil {
		return nil, err
	}

	return response, nil

}

func (is *SampleCodeServiceImpl) CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm) (*discussion_model.AiSampleCode, error) {
	// TODOC Discussion Comment 무결성 검사

	aiSampleCode, err := AiSampleCodeDbAdapter.InsertAiSampleCode(ctx, &discussion_model.CreateDiscussionAiCommentOpt{
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
func (is *SampleCodeServiceImpl) GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*AiSampleCodeResponse, error) {
	wg := new(sync.WaitGroup)

	wg.Add(AI_SAMPLE_CODE_UNIT)

	resultQueue := make(chan *AiSampleCodeResponse, AI_SAMPLE_CODE_UNIT)

	for i := 0; i < AI_SAMPLE_CODE_UNIT; i++ {
		go func(form *api.GenerateAiSampleCodesForm) {
			defer wg.Done()
			result, err := AiSampleCodeRequester.RequestReviewToAI(ctx, &AiSampleCodeRequest{
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

func (is *SampleCodeServiceImpl) DeleteAiSampleCode(ctx *context.Context, id int64) error {

	return AiSampleCodeDbAdapter.DeleteAiSampleCodeByID(ctx, id)
}
