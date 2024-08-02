package repo

import (
	"net/http"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	discussion_service "code.gitea.io/gitea/services/discussion"
	"code.gitea.io/gitea/services/forms"
)

const (
	tplDiscussionNew base.TplName = "repo/discussion/new"
)

func NewDiscussion(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true
	ctx.HTML(http.StatusOK, tplDiscussionNew)
}

func NewDiscussionPost(ctx *context.Context) {
	repo := ctx.Repo.Repository
	form := web.GetForm(ctx).(*forms.CreateDiscussionForm)

	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true
	if ctx.Written() {
		return
	}
	req := &discussion_client.PostDiscussionRequest{
		RepoId:   repo.ID,
		Poster:   ctx.Doer,
		PosterId: ctx.Doer.ID,
		Name:     form.Name,
		Content:  form.Content,
		Codes:    form.Codes,
	}
	discussion_service.NewDiscussion(ctx, repo, req)
}
