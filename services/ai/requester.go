package ai

import (
	ai_client "code.gitea.io/gitea/client/ai"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/services/context"
	"fmt"
)

const AI_SAMPLE_CODE_UNIT = 3

type AiSampleCodeRequest struct {
	CommentContent string `json:"comment"`
	CodeContent    string `json:"code"`
}

type AiSampleCodeResponse struct {
	SampleCode string `json:"sample_code"`
}

type GenerateSampleCodeResponse struct {
	SampleCode 		 string `json:"sample_code"`
	OriginalMarkdown string `json:"original_markdown"`
}

type SampleCodeRequesterImpl struct{}

type SampleCodeRequester interface {
	RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error)
}

var _ SampleCodeRequester = &SampleCodeRequesterImpl{}

type AiSampleCodeResult struct {
	SampleCode string `json:"sample_code"`
}

func (aiRequest *SampleCodeRequesterImpl) RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error) {

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

type AiPullCommentRequest struct {
	Branch   string `json:"branch"`
	TreePath string `json:"file_path"`
	Content  string `json:"code"`
}

type AiPullCommentResponse struct {
	Branch   string `json:"branch"`
	TreePath string `json:"file_path"`
	Content  string `json:"code"`
}

type PullCommentServiceImpl struct{}

type PullCommentRequesterImpl struct{}

type PullCommentRequester interface {
	RequestReviewToAI(ctx *context.Context, request *AiPullCommentRequest) (*AiPullCommentResponse, error)
}

var _ PullCommentRequester = &PullCommentRequesterImpl{}

func (aiRequest *PullCommentRequesterImpl) RequestReviewToAI(ctx *context.Context, request *AiPullCommentRequest) (*AiPullCommentResponse, error) {

	response, err := ai_client.Request().SetBody(request).Post(fmt.Sprintf("/api/pulls"))

	if err != nil {
		return nil, err
	}
	result := &AiPullCommentResponse{}

	json.Unmarshal(response.Body(), result)

	return result, nil
}
