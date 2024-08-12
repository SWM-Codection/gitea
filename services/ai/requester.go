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

type AiSampleCodeRequesterImpl struct{}

type AiSampleCodeRequester interface {
	RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error)
}

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

func (aiRequest *AiRequesterImpl) RequestReviewToAI(ctx *context.Context, request *AiReviewRequest) (*AiReviewResponse, error) {

	response, err := ai_client.Request().SetBody(request).Post(fmt.Sprintf("/api/sample"))

	if err != nil {
		return nil, err
	}
	result := &AiReviewResponse{}

	json.Unmarshal(response.Body(), result)

	return result, nil
}
