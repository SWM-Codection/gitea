package repo

// TODO 추후에 api.go로 옮기기

// import (
// 	"fmt"

// 	issues_model "code.gitea.io/gitea/models/issues"
// 	ai_service "code.gitea.io/gitea/services/ai"

// 	api "code.gitea.io/gitea/modules/structs"
// 	"code.gitea.io/gitea/modules/web"
// 	"code.gitea.io/gitea/services/context"
// )

// type AiController interface {
// 	CreateAiReviewComment(ctx *context.Context, service *ai_service.AiService, aiRequest ai_service.AiRequester, issueService ai_service.DbAdapter)
// }

// type AiControllerImpl struct{}

// func GetActionPull(ctx *context.Context) *issues_model.Issue {
// 	return GetActionIssue(ctx)
// }

// var _ AiController = &AiControllerImpl{}

// func (aiController *AiControllerImpl) CreateAiReviewComment(ctx *context.Context, service *ai_service.AiService, aiRequest ai_service.AiRequester, issueService ai_service.DbAdapter) {

// 	pull := GetActionPull(ctx)
// 	form := web.GetForm(ctx).(*api.CreateAiPullCommentForm)
// 	// TODO 결과를 받아서 AiPullComment로 저장

// 	if ctx.Written() {
// 		return
// 	}
// 	if !pull.IsPull {
// 		return
// 	}

// 	if ctx.HasError() {
// 		ctx.Flash.Error(ctx.Data["ErrorMsg"].(string))
// 		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, pull.Index))
// 		return
// 	}

// 	ai_service.AiService.CreateAiPullComment(*service, ctx, form, aiRequest, issueService)

// 	// ai_service.AiService.CreateAiPullComment(*service,
// 	// 	ctx,
// 	// 	&ai_service.AiReviewCommentResult{
// 	// 		PullID:   form.PullID,
// 	// 		RepoID:   form.RepoID,
// 	// 		Content:  result.Content,
// 	// 		TreePath: result.TreePath,
// 	// 	})

// }

// // TODOC AI 리뷰 삭제