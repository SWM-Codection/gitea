package repo

import (
	"fmt"
	"net/http"
	"strconv"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	discussion_service "code.gitea.io/gitea/services/discussion"
	"code.gitea.io/gitea/services/forms"
)

const (
	tplDiscussionNew  base.TplName = "repo/discussion/new"
	tplDiscussions    base.TplName = "repo/discussion/list"
	tplDiscussionView base.TplName = "repo/discussion/view"
)

func NewDiscussion(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true
	ctx.Data["RepoLink"] = ctx.Repo.RepoLink
	ctx.HTML(http.StatusOK, tplDiscussionNew)
}

func NewDiscussionPost(ctx *context.Context) {
	repo := ctx.Repo.Repository
	form := web.GetForm(ctx).(*forms.CreateDiscussionForm)
	log.Info("New Discussion Post Form : %v", form)
	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true
	if ctx.Written() {
		return
	}
	req := &discussion_client.PostDiscussionRequest{
		RepoId:     repo.ID,
		Poster:     ctx.Doer,
		PosterId:   ctx.Doer.ID,
		Name:       form.Name,
		Content:    form.Content,
		BranchName: form.BranchName,
		Codes:      form.Codes,
	}
	discussionId, err := discussion_service.NewDiscussion(ctx, repo, req)
	log.Info("New Discussion Post : %v", discussionId)
	if err != nil {
		ctx.ServerError("NewDiscussion", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"discussionId": discussionId,
	})
}

// TODO: for now, some clumsy logics are included, but for later this function should be polished
func Discussions(ctx *context.Context) {
	// get discussion lists
	listResp, err := discussion_service.GetDiscussionList(ctx)
	if err != nil {
		ctx.ServerError("Discussions", err)
		return
	}

	// make paginator
	currentPage, err := strconv.Atoi(ctx.FormString("page"))
	if err != nil {
		currentPage = 1
	}
	pager := context.NewPagination(int(listResp.TotalCount), setting.UI.IssuePagingNum, currentPage, 5)

	// get count info
	count, err := discussion_client.GetDiscussionCount(ctx.Repo.Repository.ID)
	if err != nil {
		ctx.ServerError("Discussions:GetDiscussionCount", err)
	}

	// get state info
	isClosed := ctx.FormString("state") == "closed"
	state := "open"
	if isClosed {
		state = "closed"
	}

	// prepare data for tamplete
	ctx.Data["Title"] = ctx.Tr("repo.discussion.list")
	ctx.Data["PageIsDiscussionList"] = true
	ctx.Data["Discussions"] = listResp.Discussions
	ctx.Data["RepoLink"] = ctx.Repo.RepoLink
	ctx.Data["Page"] = pager
	ctx.Data["OpenCount"] = count.OpenCount
	ctx.Data["ClosedCount"] = count.CloseCount
	ctx.Data["State"] = state
	ctx.Data["AllStatesLink"] = fmt.Sprintf("%s/discussions?state=%s&page=1", ctx.Repo.RepoLink, state)
	ctx.Data["OpenLink"] = fmt.Sprintf("%s/discussions?state=open&page=1", ctx.Repo.RepoLink)
	ctx.Data["ClosedLink"] = fmt.Sprintf("%s/discussions?state=closed&page=1", ctx.Repo.RepoLink)

	ctx.HTML(http.StatusOK, tplDiscussions)
}

func ViewDiscussion(ctx *context.Context) {
	discussionId := ctx.ParamsInt64(":index")
	discussionResponse, err := discussion_client.GetDiscussion(discussionId)
	log.Info("discussion: %v, err: %v", discussionResponse, err)
	ctx.Data["Discussion"] = discussionResponse
	ctx.HTML(http.StatusOK, tplDiscussionView)
}
