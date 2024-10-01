// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	api "code.gitea.io/gitea/modules/structs"
	issues_model "code.gitea.io/gitea/models/issues"
	ai_service "code.gitea.io/gitea/services/ai"
	pull_model "code.gitea.io/gitea/models/pull"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/json"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/context/upload"
	"code.gitea.io/gitea/services/forms"
	"code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/modules/git"
	pull_service "code.gitea.io/gitea/services/pull"
	user_service "code.gitea.io/gitea/services/user"
)

const (
	tplDiffConversation     base.TplName = "repo/diff/conversation"
	tplConversationOutdated base.TplName = "repo/diff/conversation_outdated"
	tplTimelineConversation base.TplName = "repo/issue/view_content/conversation"
	tplNewComment           base.TplName = "repo/diff/new_comment"
)

// RenderNewCodeCommentForm will render the form for creating a new review comment
func RenderNewCodeCommentForm(ctx *context.Context) {
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}
	if !issue.IsPull {
		return
	}
	currentReview, err := issues_model.GetCurrentReview(ctx, ctx.Doer, issue)
	if err != nil && !issues_model.IsErrReviewNotExist(err) {
		ctx.ServerError("GetCurrentReview", err)
		return
	}
	ctx.Data["PageIsPullFiles"] = true
	ctx.Data["Issue"] = issue
	ctx.Data["CurrentReview"] = currentReview
	pullHeadCommitID, err := ctx.Repo.GitRepo.GetRefCommitID(issue.PullRequest.GetGitRefName())
	if err != nil {
		ctx.ServerError("GetRefCommitID", err)
		return
	}
	ctx.Data["AfterCommitID"] = pullHeadCommitID
	ctx.Data["IsAttachmentEnabled"] = setting.Attachment.Enabled
	upload.AddUploadContext(ctx, "comment")
	ctx.HTML(http.StatusOK, tplNewComment)
}

// 체크 코드 코멘트를 만드는 곳
// CreateCodeComment will create a code comment including an pending review if required
func CreateCodeComment(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.CodeCommentForm)
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}
	if !issue.IsPull {
		return
	}

	if ctx.HasError() {
		ctx.Flash.Error(ctx.Data["ErrorMsg"].(string))
		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		return
	}

	signedLine := form.Line
	if form.Side == "previous" {
		signedLine *= -1
	}

	var attachments []string
	if setting.Attachment.Enabled {
		attachments = form.Files
	}

	comment, err := pull_service.CreateCodeComment(ctx,
		ctx.Doer,
		ctx.Repo.GitRepo,
		issue,
		signedLine,
		form.Content,
		form.TreePath,
		!form.SingleReview,
		form.Reply,
		form.LatestCommitID,
		attachments,
	)
	if err != nil {
		ctx.ServerError("CreateCodeComment", err)
		return
	}

	if comment == nil {
		log.Trace("Comment not created: %-v #%d[%d]", ctx.Repo.Repository, issue.Index, issue.ID)
		ctx.Redirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		return
	}

	log.Trace("Comment created: %-v #%d[%d] Comment[%d]", ctx.Repo.Repository, issue.Index, issue.ID, comment.ID)

	renderConversation(ctx, comment, form.Origin)
}

// UpdateResolveConversation add or remove an Conversation resolved mark
func UpdateResolveConversation(ctx *context.Context) {
	origin := ctx.FormString("origin")
	action := ctx.FormString("action")
	commentID := ctx.FormInt64("comment_id")

	comment, err := issues_model.GetCommentByID(ctx, commentID)
	if err != nil {
		ctx.ServerError("GetIssueByID", err)
		return
	}

	if err = comment.LoadIssue(ctx); err != nil {
		ctx.ServerError("comment.LoadIssue", err)
		return
	}

	if comment.Issue.RepoID != ctx.Repo.Repository.ID {
		ctx.NotFound("comment's repoID is incorrect", errors.New("comment's repoID is incorrect"))
		return
	}

	var permResult bool
	if permResult, err = issues_model.CanMarkConversation(ctx, comment.Issue, ctx.Doer); err != nil {
		ctx.ServerError("CanMarkConversation", err)
		return
	}
	if !permResult {
		ctx.Error(http.StatusForbidden)
		return
	}

	if !comment.Issue.IsPull {
		ctx.Error(http.StatusBadRequest)
		return
	}

	if action == "Resolve" || action == "UnResolve" {
		err = issues_model.MarkConversation(ctx, comment, ctx.Doer, action == "Resolve")
		if err != nil {
			ctx.ServerError("MarkConversation", err)
			return
		}
	} else {
		ctx.Error(http.StatusBadRequest)
		return
	}

	renderConversation(ctx, comment, origin)
}

func CreateAiPullSampleCode(ctx *context.Context) {
	// TODOC swagger 추가
	// TODOC 공격 우려가 있어서 Create할 비대칭키 방식 암호화가 필요해보임.
	form := web.GetForm(ctx).(*api.CreateAiSampleCodesForm)

	targetCommentId, err := strconv.ParseInt(form.TargetCommentId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": "Invalid TargetCommentId format",
		})
		return
	}

	sampleCode, err := ai_service.AiSampleCodeService.CreateAiSampleCode(ctx, form)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err.Error(),
		})
		return
	}
	println(sampleCode.Id)

	targetComment, err := issues_model.GetCommentByID(ctx, targetCommentId)
	if err != nil {
		ctx.ServerError("GetIssueByID", err)
		return
	}

	if err = targetComment.LoadIssue(ctx); err != nil {
		ctx.ServerError("comment.LoadIssue", err)
		return
	}

	if err = targetComment.Issue.LoadRepo(ctx); err != nil {
		ctx.ServerError("Issue.LoadRepo", err)
		return
	}
	
	initializeRepo(ctx, targetComment.Issue.Repo.OwnerName, targetComment.Issue.Repo.Name)
	renderConversation(ctx, targetComment, form.OriginData)
}

// 체크 코멘트가 만들어졌을 때 렌더링 되는 곳
func renderConversation(ctx *context.Context, comment *issues_model.Comment, origin string) {
	ctx.Data["PageIsPullFiles"] = origin == "diff"

	showOutdatedComments := origin == "timeline" || ctx.Data["ShowOutdatedComments"].(bool)
	comments, err := issues_model.FetchCodeCommentsByLine(ctx, comment.Issue, ctx.Doer, comment.TreePath, comment.Line, showOutdatedComments)
	if err != nil {
		ctx.ServerError("FetchCodeCommentsByLine", err)
		return
	}
	aiPullComment, err := issues_model.FetchAiPullCommentByLine(ctx, comment.Issue, comment.TreePath, comment.Line)

	if err != nil {
		ctx.ServerError("FetchAiPullCommentByLine", err)
		return
	}
	if aiPullComment != nil {
		if err := aiPullComment.LoadPoster(ctx); err != nil {
			ctx.ServerError("LoadPoster", err)
			return
		}
		if err := aiPullComment.LoadAttachments(ctx); err != nil {
			ctx.ServerError("LoadPoster", err)
			return
		}
		println(comments[0].RenderedContent)
		println()
		comments = append([]*issues_model.Comment{aiPullComment}, comments...)
		println(comments[0].RenderedContent)
	}

	if len(comments) == 0 {
		// if the comments are empty (deleted, outdated, etc), it's better to tell the users that it is outdated
		ctx.HTML(http.StatusOK, tplConversationOutdated)
		return
	}

	if err := comments.LoadAttachments(ctx); err != nil {
		ctx.ServerError("LoadAttachments", err)
		return
	}

	ctx.Data["IsAttachmentEnabled"] = setting.Attachment.Enabled
	upload.AddUploadContext(ctx, "comment")

	ctx.Data["comments"] = comments
	if ctx.Data["CanMarkConversation"], err = issues_model.CanMarkConversation(ctx, comment.Issue, ctx.Doer); err != nil {
		ctx.ServerError("CanMarkConversation", err)
		return
	}
	ctx.Data["Issue"] = comment.Issue
	if err = comment.Issue.LoadPullRequest(ctx); err != nil {
		ctx.ServerError("comment.Issue.LoadPullRequest", err)
		return
	}
	pullHeadCommitID, err := ctx.Repo.GitRepo.GetRefCommitID(comment.Issue.PullRequest.GetGitRefName())
	if err != nil {
		ctx.ServerError("GetRefCommitID", err)
		return
	}
	ctx.Data["AfterCommitID"] = pullHeadCommitID
	ctx.Data["CanBlockUser"] = func(blocker, blockee *user_model.User) bool {
		return user_service.CanBlockUser(ctx, ctx.Doer, blocker, blockee)
	}

	if origin == "diff" {
		ctx.HTML(http.StatusOK, tplDiffConversation)
	} else if origin == "timeline" {
		ctx.HTML(http.StatusOK, tplTimelineConversation)
	} else {
		ctx.Error(http.StatusBadRequest, "Unknown origin: "+origin)
	}
}

// SubmitReview creates a review out of the existing pending review or creates a new one if no pending review exist
func SubmitReview(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.SubmitReviewForm)
	issue := GetActionIssue(ctx)
	if ctx.Written() {
		return
	}
	if !issue.IsPull {
		return
	}
	if ctx.HasError() {
		ctx.Flash.Error(ctx.Data["ErrorMsg"].(string))
		ctx.JSONRedirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		return
	}

	reviewType := form.ReviewType()
	switch reviewType {
	case issues_model.ReviewTypeUnknown:
		ctx.ServerError("ReviewType", fmt.Errorf("unknown ReviewType: %s", form.Type))
		return

	// can not approve/reject your own PR
	case issues_model.ReviewTypeApprove, issues_model.ReviewTypeReject:
		if issue.IsPoster(ctx.Doer.ID) {
			var translated string
			if reviewType == issues_model.ReviewTypeApprove {
				translated = ctx.Locale.TrString("repo.issues.review.self.approval")
			} else {
				translated = ctx.Locale.TrString("repo.issues.review.self.rejection")
			}

			ctx.Flash.Error(translated)
			ctx.JSONRedirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
			return
		}
	}

	var attachments []string
	if setting.Attachment.Enabled {
		attachments = form.Files
	}

	_, comm, err := pull_service.SubmitReview(ctx, ctx.Doer, ctx.Repo.GitRepo, issue, reviewType, form.Content, form.CommitID, attachments)
	if err != nil {
		if issues_model.IsContentEmptyErr(err) {
			ctx.Flash.Error(ctx.Tr("repo.issues.review.content.empty"))
			ctx.JSONRedirect(fmt.Sprintf("%s/pulls/%d/files", ctx.Repo.RepoLink, issue.Index))
		} else if errors.Is(err, pull_service.ErrSubmitReviewOnClosedPR) {
			ctx.Status(http.StatusUnprocessableEntity)
		} else {
			ctx.ServerError("SubmitReview", err)
		}
		return
	}
	ctx.JSONRedirect(fmt.Sprintf("%s/pulls/%d#%s", ctx.Repo.RepoLink, issue.Index, comm.HashTag()))
}

// DismissReview dismissing stale review by repo admin
func DismissReview(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.DismissReviewForm)
	comm, err := pull_service.DismissReview(ctx, form.ReviewID, ctx.Repo.Repository.ID, form.Message, ctx.Doer, true, true)
	if err != nil {
		if pull_service.IsErrDismissRequestOnClosedPR(err) {
			ctx.Status(http.StatusForbidden)
			return
		}
		ctx.ServerError("pull_service.DismissReview", err)
		return
	}

	ctx.Redirect(fmt.Sprintf("%s/pulls/%d#%s", ctx.Repo.RepoLink, comm.Issue.Index, comm.HashTag()))
}

// viewedFilesUpdate Struct to parse the body of a request to update the reviewed files of a PR
// If you want to implement an API to update the review, simply move this struct into modules.
type viewedFilesUpdate struct {
	Files         map[string]bool `json:"files"`
	HeadCommitSHA string          `json:"headCommitSHA"`
}

func UpdateViewedFiles(ctx *context.Context) {
	// Find corresponding PR
	issue, ok := getPullInfo(ctx)
	if !ok {
		return
	}
	pull := issue.PullRequest

	var data *viewedFilesUpdate
	err := json.NewDecoder(ctx.Req.Body).Decode(&data)
	if err != nil {
		log.Warn("Attempted to update a review but could not parse request body: %v", err)
		ctx.Resp.WriteHeader(http.StatusBadRequest)
		return
	}

	// Expect the review to have been now if no head commit was supplied
	if data.HeadCommitSHA == "" {
		data.HeadCommitSHA = pull.HeadCommitID
	}

	updatedFiles := make(map[string]pull_model.ViewedState, len(data.Files))
	for file, viewed := range data.Files {
		// Only unviewed and viewed are possible, has-changed can not be set from the outside
		state := pull_model.Unviewed
		if viewed {
			state = pull_model.Viewed
		}
		updatedFiles[file] = state
	}

	if err := pull_model.UpdateReviewState(ctx, ctx.Doer.ID, pull.ID, data.HeadCommitSHA, updatedFiles); err != nil {
		ctx.ServerError("UpdateReview", err)
	}
}

func initializeRepo(ctx *context.Context, ownerName string, repoName string) {
    if ctx.Repo.Repository == nil {
        repo, err := repo.GetRepositoryByOwnerAndName(ctx, ownerName, repoName)
        if err != nil {
            ctx.ServerError("GetRepositoryByOwnerAndName", err)
            return
        }
        ctx.Repo.Repository = repo
    }

    if ctx.Repo.GitRepo == nil {
        gitRepo, err := git.OpenRepository(ctx, ctx.Repo.Repository.RepoPath())
        if err != nil {
            ctx.ServerError("OpenRepository", err)
            return
        }
        ctx.Repo.GitRepo = gitRepo
		defer ctx.Repo.GitRepo.Close()
    }
}