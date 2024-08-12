package repo

import (
	"github.com/spf13/cast"
	"net/http"

	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/web"
	ai_service "code.gitea.io/gitea/services/ai"

	"code.gitea.io/gitea/services/context"
)

// TODOC 특정한 PR에 대한 AI comment 리퀘스트가 이미 이루어졌을 경우를 체크해서 중복 요청 차단.
func CreateAiPullComment(ctx *context.Context) {
	// swagger:operation POST /ai/pull/review repository repoCreateAiPullComment
	// ---
	// summary: Create ai pull comment
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/CreateAiPullCommentForm"
	// responses:
	//   "201":
	//     "$ref": "#/responses/Attachment"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"

	// Check if attachments are enabled

	form := web.GetForm(ctx).(*api.CreateAiPullCommentForm)

	// TODOC 싱글톤으로 바꾼 뒤에
	// TODOC 설정을 통해 의존성 주입하는 방식으로 바꾸기

	aiService := new(ai_service.AiServiceImpl)
	aiRequester := new(ai_service.AiRequesterImpl)
	adapter := new(ai_service.DbAdapterImpl)
	err := ai_service.AiService.CreateAiPullComment(aiService, ctx, form, aiRequester, adapter)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, map[string]any{
		"message": "request has accepted",
	})
}

// TODOC Delete 엔드포인트 작성

func GenerateAiSampleCodes(ctx *context.Context) {
	// swagger:operation POST /ai/pull/review repository repoCreateAiPullComment
	// ---
	// summary: Create ai pull comment
	// produces:
	// - application/json
	// consumes:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/CreateAiPullCommentForm"
	// responses:
	//   "201":
	//     "$ref": "#/responses/Attachment"
	//   "400":
	//     "$ref": "#/responses/error"
	//   "404":
	//     "$ref": "#/responses/notFound"

	// Check if attachments are enabled

	form := web.GetForm(ctx).(*api.GenerateAiSampleCodesForm)

	// TODOC 싱글톤으로 바꾼 뒤에
	// TODOC 설정을 통해 의존성 주입하는 방식으로 바꾸기

	aiService := new(ai_service.DiscussionAiServiceImpl)
	aiRequester := new(ai_service.AiSampleCodeRequesterImpl)
	adapter := new(ai_service.DiscussionDbAdapterImpl)
	sampleCodes, err := ai_service.DiscussionAiService.GenerateAiSampleCodes(aiService, ctx, form, aiRequester, adapter)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error()},
		)
		return
	}

	ctx.JSON(http.StatusAccepted, sampleCodes)
}

func CreateAiSampleCode(ctx *context.Context) {

	// TODOC swagger 추가
	// TODOC 공격 우려가 있어서 Create할 비대칭키 방식 암호화가 필요해보임.
	form := web.GetForm(ctx).(*api.CreateAiSampleCodesForm)

	aiService := new(ai_service.DiscussionAiServiceImpl)
	adapter := new(ai_service.DiscussionDbAdapterImpl)
	sampleCode, err := ai_service.DiscussionAiService.CreateAiSampleCode(aiService, ctx, form, adapter)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusAccepted, sampleCode)

}

func DeleteAiSampleCode(ctx *context.Context) {

	form := web.GetForm(ctx).(*api.DeleteSampleCodesForm)

	aiService := new(ai_service.DiscussionAiServiceImpl)
	adapter := new(ai_service.DiscussionDbAdapterImpl)

	err := ai_service.DiscussionAiService.DeleteAiSampleCode(aiService, ctx, cast.ToInt64(form.TargetCommentId), adapter)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusAccepted, map[string]any{})
}
