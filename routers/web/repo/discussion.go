package repo

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	ai_service "code.gitea.io/gitea/services/ai"

	"code.gitea.io/gitea/models/discussion"
	"gitea.com/go-chi/binding"

	stdCtx "context"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/client/discussion/model"
	issues_model "code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/models/organization"
	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/markup/markdown"
	"code.gitea.io/gitea/modules/optional"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	discussion_service "code.gitea.io/gitea/services/discussion"
	"code.gitea.io/gitea/services/forms"
	notify_service "code.gitea.io/gitea/services/notify"
)

const (
	tplDiscussionNew            base.TplName = "repo/discussion/new"
	tplDiscussions              base.TplName = "repo/discussion/list"
	tplDiscussionView           base.TplName = "repo/discussion/view"
	tplDiscussionFileComments   base.TplName = "repo/discussion/file_comments"
	tplDiscussionFiles          base.TplName = "repo/discussion/view_file"
	tplNewDiscussionFileComment base.TplName = "repo/discussion/new_file_comment"
	tplDiscussionFileComment    base.TplName = "repo/discussion/file_comment_holder"
	tplDiscussionReactions      base.TplName = "repo/discussion/reactions"
	tplAiCommentForm            base.TplName = "repo/discussion/new_ai_comment"
)

func GetAiDiscussionForm(ctx *context.Context) {

	queryParams := ctx.Req.URL.Query()

	discussionId := queryParams.Get("discussionId")
	codeId := queryParams.Get("codeId")
	startLine := queryParams.Get("startLine")
	endLine := queryParams.Get("endLine")

	ctx.Data["DiscussionId"] = discussionId
	ctx.Data["CodeId"] = codeId
	ctx.Data["StartLine"] = startLine
	ctx.Data["EndLine"] = endLine
	ctx.PageData["RepoLink"] = ctx.Repo.RepoLink
	ctx.Data["Repository"] = ctx.Repo.Repository
	ctx.HTML(http.StatusOK, tplAiCommentForm)
}

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
	req := &model.PostDiscussionRequest{
		RepoId:     repo.ID,
		Poster:     ctx.Doer,
		PosterId:   ctx.Doer.ID,
		Name:       form.Name,
		Content:    form.Content,
		BranchName: form.BranchName,
		Codes:      form.Codes,
	}
	discussionId, err := discussion_service.NewDiscussion(ctx, repo, req)
	if err != nil {
		ctx.ServerError("NewDiscussion", err)
		return
	}
	notify_service.NewDiscussion(ctx, ctx.Doer, repo, discussionId)
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
	pinned, err := discussion_service.GetPinnedDiscussionList(ctx)
	if err != nil {
		ctx.ServerError("Discussions:GetPinnedDiscussions", err)
	}
	var isShowClosed optional.Option[bool]
	switch ctx.FormString("state") {
	case "closed":
		isShowClosed = optional.Some(true)
	case "all":
		isShowClosed = optional.None[bool]()
	default:
		isShowClosed = optional.Some(false)
	}
	// if there are closed discussions and no open discussionss, default to showing all discusisons
	if len(ctx.FormString("state")) == 0 && listResp.TotalCount == 0 {
		isShowClosed = optional.None[bool]()
	}

	log.Info("discussions : %v", listResp.Discussions)
	// prepare data for tamplete
	ctx.Data["IsRepoAdmin"] = ctx.IsSigned && (ctx.Repo.IsAdmin() || ctx.Doer.IsAdmin)
	ctx.Data["PinnedDiscussions"] = pinned.Discussions
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
	ctx.Data["IsShowClosed"] = isShowClosed
	switch {
	case isShowClosed.Value():
		ctx.Data["State"] = "closed"
	case !isShowClosed.Has():
		ctx.Data["State"] = "all"
	default:
		ctx.Data["State"] = "open"
	}
	ctx.HTML(http.StatusOK, tplDiscussions)
}

func ViewDiscussion(ctx *context.Context) {
	discussionId := ctx.ParamsInt64(":index")
	var assignees = make([]*user_model.User, 0, 10)
	var participants = make([]*user_model.User, 1, 10)

	discussionResponse, err := discussion_client.GetDiscussion(discussionId)
	if err != nil {
		ctx.ServerError("error on discussion response: err = %v", err)
	}
	// valid dicussion id must bigger than 0
	if discussionResponse.Id == 0 {
		ctx.NotFound("discussion not exists", fmt.Errorf("discussion with %v does not exist", discussionResponse.Id))
		return
	}
	poster, err := user_model.GetUserByID(ctx, discussionResponse.PosterId)
	if err != nil {
		ctx.ServerError("errro on get user by id: err = %v", err)
	}
	discussionResponse.Poster = poster
	discussionContentResponse, err := discussion_client.GetDiscussionContents(discussionId)
	if err != nil {
		ctx.ServerError("error on discussion content response: err = %v", err)
		return
	}
	rd, err := discussionRoleDescriptor(ctx, ctx.Repo.Repository, discussionResponse.Poster, discussionResponse)
	if err != nil {
		ctx.ServerError("error on retreiving discussion role descriptor, %v", err)
		return
	}
	for _, assigneeId := range discussionResponse.Assignees {
		assignee, err := user_model.GetUserByID(ctx, assigneeId)
		if err != nil {
			ctx.ServerError("errro on get user by id: err = %v", err)
		}
		assignees = append(assignees, assignee)
	}
	repo := ctx.Repo.Repository
	assigneeUsers, err := repo_model.GetRepoAssignees(ctx, repo)
	if err != nil {
		ctx.ServerError("GetRepoAssignees", err)
		return
	}

	var pinAllowed bool
	pinAllowed, err = discussion_client.IsNewPinAllowed(discussionResponse.RepoId)
	if err != nil {
		ctx.ServerError("IsNewPinAllowed", err)
		return
	}

	participants[0] = poster
	ctx.Data["IsRepoAdmin"] = ctx.IsSigned && (ctx.Repo.IsAdmin() || ctx.Doer.IsAdmin)
	ctx.Data["DiscussionContent"] = discussionContentResponse
	ctx.Data["PageIsDiscussionList"] = true
	ctx.Data["Repository"] = ctx.Repo.Repository
	ctx.Data["Discussion"] = discussionResponse
	ctx.Data["DiscussionTab"] = "conversation"
	ctx.Data["DiscussionRoleDescriptor"] = rd
	ctx.Data["DiscussionAssignees"] = assignees
	ctx.Data["Participants"] = participants
	ctx.Data["NumParticipants"] = len(participants)
	ctx.Data["Assignees"] = MakeSelfOnTop(ctx.Doer, assigneeUsers)
	ctx.Data["NewPinAllowed"] = pinAllowed
	ctx.HTML(http.StatusOK, tplDiscussionView)
}

func ViewDiscussionFiles(ctx *context.Context) {

	discussionId := ctx.ParamsInt64(":index")
	discussionResponse, err := discussion_client.GetDiscussion(discussionId)
	if err != nil {
		ctx.ServerError("error on discussion response: err = %v", err)
	}
	poster, err := user_model.GetUserByID(ctx, discussionResponse.PosterId)
	if err != nil {
		ctx.ServerError("errro on get user by id: err = %v", err)
	}
	discussionResponse.Poster = poster

	ctx.Data["PageIsDiscussionList"] = true
	ctx.Data["Repository"] = ctx.Repo.Repository
	ctx.Data["Discussion"] = discussionResponse
	ctx.PageData["RepoLink"] = ctx.Repo.RepoLink
	ctx.PageData["DiscussionId"] = discussionId
	ctx.Data["DiscussionTab"] = "files"

	ctx.HTML(http.StatusOK, tplDiscussionView)
}

func NewDiscussionCommentPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.CreateDiscussionCommentForm)

	discussionId := ctx.ParamsInt64(":discussionId")
	commentScope := util.Iif(
		form.CodeId != nil && form.StartLine != nil && form.EndLine != nil,
		model.CommentScopeLocal,
		model.CommentScopeGlobal,
	)
	// if request malformed, then make it valid
	if commentScope == model.CommentScopeGlobal {
		form.CodeId = nil
		form.StartLine = nil
		form.EndLine = nil
		form.GroupId = nil
	}

	req := &model.PostCommentRequest{
		DiscussionId: discussionId,
		Scope:        commentScope,
		PosterId:     ctx.Doer.ID,
		Content:      form.Content,
		GroupId:      form.GroupId,
		CodeId:       form.CodeId,
		StartLine:    form.StartLine,
		EndLine:      form.EndLine,
	}
	id, err := discussion_client.PostComment(req)
	if err != nil {
		ctx.JSONError(fmt.Sprintf("failed to post discussion comment %v", err))
	}
	if id == nil {
		// XXX check reachability later
		ctx.JSONError("hmm something weird..")
		return
	}
	ctx.JSON(http.StatusOK, map[string]int64{"id": *id})

}

// TODO: 추후 NewDiscussionPost 메소드와 통합할지 고려해보기
func RenderNewDiscussionComment(ctx *context.Context) {

	id := ctx.ParamsInt64("id")
	comment, err := discussion_client.GetDiscussionComment(id)

	if err != nil {
		ctx.ServerError("failed to fetch comment: %v", err)
		return
	}

	poster, err := user_model.GetUserByID(ctx, comment.PosterId)

	if err != nil {
		ctx.ServerError("failed to fetch user data: %v", err)
		return
	}

	// TODO: 답글 기능 고려해서 넣기
	comments := make([]*discussion_service.DiscussionComment, 0, 1)

	newComment := &discussion_service.DiscussionComment{
		ID:           comment.Id,
		StartLine:    comment.StartLine,
		DiscussionId: comment.DiscussionId,
		GroupId:      comment.GroupId,
		EndLine:      comment.EndLine,
		CodeId:       comment.CodeId,
		CreatedUnix:  comment.CreatedUnix,
		Reactions:    comment.Reactions,
		Poster:       poster,
		Content:      comment.Content,
	}
	newComment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
		Ctx: ctx,
		Links: markup.Links{
			Base: ctx.Repo.RepoLink,
		},
	}, newComment.Content)

	if err != nil {
		ctx.ServerError("markdown rendering failed : %v", err)
	}

	comments = append(comments, newComment)
	ctx.Data["comments"] = comments
	ctx.Data["DiscussionId"] = newComment.DiscussionId
	if newComment.GroupId == newComment.ID {
		ctx.HTML(http.StatusOK, tplDiscussionFileComments)
	} else {
		ctx.HTML(http.StatusOK, tplDiscussionFileComment)
	}
}

type DiscussionFileCommentsResponse struct {
	Html    *template.HTML `json:"html"`
	EndLine int64          `json:"endLine"`
}

func groupDiscussionFileCommentsByGroupId(
	ctx *context.Context,
	commentsResp []*discussion_service.DiscussionComment) (map[int64][]*discussion_service.DiscussionComment, error) {

	comments := make(map[int64][]*discussion_service.DiscussionComment)

	for _, comment := range commentsResp {
		var err error
		comment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: ctx.Repo.RepoLink,
			},
		}, comment.Content)

		if err != nil {
			return nil, err
		}

		comments[comment.GroupId] = append(comments[comment.GroupId], comment)

	}

	return comments, nil
}

func renderDiscussionFileComments(ctx *context.Context, commentGroups map[int64][]*discussion_service.DiscussionComment) ([]*DiscussionFileCommentsResponse, error) {

	resp := make([]*DiscussionFileCommentsResponse, 0, len(commentGroups))

	for _, groupComments := range commentGroups {

		sort.Slice(groupComments, func(i, j int) bool {
			return groupComments[i].CreatedUnix < groupComments[j].CreatedUnix
		})
		ctx.Data["RepoLink"] = ctx.Repo.RepoLink
		ctx.Data["comments"] = groupComments
		ctx.Data["DiscussionId"] = groupComments[0].DiscussionId
		html, err := ctx.RenderToHTML(tplDiscussionFileComments, ctx.Data)

		if err != nil {
			return nil, err
		}

		endLine := groupComments[0].EndLine

		resp = append(resp, &DiscussionFileCommentsResponse{
			Html:    &html,
			EndLine: endLine,
		})

	}

	return resp, nil

}

func DiscussionComments(ctx *context.Context) {
	codeId := ctx.ParamsInt64("codeId")

	comments, err := discussion_service.GetDiscussionCommentsByCodeId(ctx, codeId)

	if err != nil {
		ctx.JSONError(err)
	}

	commentGroups, err := groupDiscussionFileCommentsByGroupId(ctx, comments)

	if err != nil {
		ctx.JSONError(err)
	}

	resp, err := renderDiscussionFileComments(ctx, commentGroups)

	if err != nil {
		ctx.JSONError(err)
	}

	ctx.JSON(http.StatusOK, resp)

}

func RenderNewDiscussionFileCommentForm(ctx *context.Context) {
	// RenderNewCodeCommentForm will render the form for creating a new review comment

	queryParams := ctx.Req.URL.Query()

	discussionId := queryParams.Get("discussionId")
	codeId := queryParams.Get("codeId")
	startLine := queryParams.Get("startLine")
	endLine := queryParams.Get("endLine")

	ctx.Data["DiscussionId"] = discussionId
	ctx.Data["CodeId"] = codeId
	ctx.Data["StartLine"] = startLine
	ctx.Data["EndLine"] = endLine
	ctx.PageData["RepoLink"] = ctx.Repo.RepoLink
	ctx.Data["Repository"] = ctx.Repo.Repository
	ctx.HTML(http.StatusOK, tplNewDiscussionFileComment)
}

func DiscussionContent(ctx *context.Context) {

	discussionId := ctx.ParamsInt64(":index")

	discussionContent, err := discussion_service.GetDiscussionContentWithHighlights(discussionId)

	if err != nil {
		ctx.JSONError(err.Error())
	}

	ctx.JSON(http.StatusOK, discussionContent)
}

func SetDiscussionClosedState(ctx *context.Context) {
	repo := ctx.Repo.Repository
	discussionId := ctx.ParamsInt64(":discussionId")
	queryParams := ctx.Req.URL.Query()
	isClosedStr := queryParams.Get("isClosed")

	isClosed, err := strconv.ParseBool(isClosedStr)
	if err != nil {
		ctx.ServerError("Invalid 'isClosed' parameter", err)
		return
	}

	err = discussion_client.SetDiscussionClosedState(discussionId, isClosed)
	if err != nil {
		ctx.ServerError(fmt.Sprintf("Failed to set review state for discussion %d", discussionId), err)
		return
	}

	// update close discussion numbers
	if err := repo_model.UpdateRepoDiscussionNumbers(ctx, repo.ID, true); err != nil {
		ctx.ServerError("Failed to update review state for close discussion", err)
		return
	}

	ctx.Status(http.StatusOK)
}

func SetDiscussionDeadline(ctx *context.Context) {
	form := web.GetForm(ctx).(*api.EditDeadlineOption)
	discussionId := ctx.ParamsInt64(":discussionId")
	var deadlineUnix int64 // Unix 타임스탬프를 저장할 int64 변수
	var deadline time.Time
	if form.Deadline != nil && !form.Deadline.IsZero() {
		deadline = time.Date(form.Deadline.Year(), form.Deadline.Month(), form.Deadline.Day(),
			23, 59, 59, 0, time.Local)
		deadlineUnix = deadline.Unix() // Unix 타임스탬프를 int64로 저장
	}

	err := discussion_client.SetDiscussionDeadline(discussionId, deadlineUnix)
	if err != nil {
		ctx.ServerError(fmt.Sprintf("Failed to set review state for discussion %d", discussionId), err)
		return
	}

	ctx.JSON(http.StatusCreated, model.DiscussionDeadline{Deadline: &deadline})
}

func getActionDiscussionIds(ctx *context.Context) []int64 {
	commaSeparatedDiscussionIDs := ctx.FormString("issue_ids")
	if len(commaSeparatedDiscussionIDs) == 0 {
		return nil
	}
	discussionIDs := make([]int64, 0, 10)
	for _, stringDiscussionID := range strings.Split(commaSeparatedDiscussionIDs, ",") {
		discussionID, err := strconv.ParseInt(stringDiscussionID, 10, 64)
		if err != nil {
			ctx.ServerError("ParseInt", err)
			return nil
		}
		discussionIDs = append(discussionIDs, discussionID)
	}
	return discussionIDs
}

func UpdateDiscussionStatus(ctx *context.Context) {
	discussionIds := getActionDiscussionIds(ctx)
	repo := ctx.Repo.Repository

	var isClosed bool
	switch action := ctx.FormString("action"); action {
	case "open":
		isClosed = false
	case "close":
		isClosed = true
	default:
		log.Warn("Unrecognized action: %s", action)
	}

	for _, discussionId := range discussionIds {
		err := discussion_client.SetDiscussionClosedState(discussionId, isClosed)
		if err != nil {
			ctx.ServerError(fmt.Sprintf("Failed to set review state for discussion %d", discussionId), err)
			return
		}
	}

	// update close discussion numbers
	if err := repo_model.UpdateRepoDiscussionNumbers(ctx, repo.ID, true); err != nil {
		ctx.ServerError("Failed to update review state for close discussion", err)
		return
	}

	ctx.JSONOK()
}

func DeleteDiscussionFileComment(ctx *context.Context) {
	posterId := ctx.Doer.ID

	discussionCommentId := ctx.ParamsInt64(":discussionId")

	if discussionCommentId < 0 {
		err := discussion.DeleteAiSampleCodeByID(ctx, -discussionCommentId)

		if err != nil {
			ctx.JSONError(err.Error())
			return
		}
	} else {
		err := discussion_service.DeleteDiscussionComment(ctx, discussionCommentId, posterId)

		if err != nil {
			ctx.JSONError(err.Error())
			return
		}
	}

	ctx.JSONOK()

}

func ModifyDiscussionFileComment(ctx *context.Context) {

	form := web.GetForm(ctx).(*forms.ModifyDiscussionCommentForm)

	formErr := form.Validate(ctx.Req, binding.Errors{})

	if formErr.Len() > 0 {
		log.Error("ModifyDiscussionFileComment: %v", formErr)
		ctx.JSONErrorf("ModifyDiscussionFileComment: %v", formErr)
	}

	var err error

	if form.DiscussionCommentId < 0 {
		err = ai_service.UpdateAiSampleCode(ctx, form)

	} else {
		err = discussion_service.ModifyDiscussionComment(ctx, form)

	}

	if err != nil {
		ctx.ServerError("modify request failed %v", err)
		return
	}

	var renderedContent template.HTML
	if form.Content != "" {
		renderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: ctx.Repo.RepoLink,
			},
		}, form.Content)
		if err != nil {
			ctx.ServerError("RenderString", err)
			return
		}

	} else {
		contentEmpty := fmt.Sprintf(`<span class="no-content">%s</span>`, ctx.Tr("repo.issues.no_content"))
		renderedContent = template.HTML(contentEmpty)
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"content": renderedContent,
	})

}

func ChangeDiscussionCommentReaction(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.ReactionForm)
	if ctx.Written() {
		return
	}
	discussionId, err := strconv.ParseInt(ctx.Params(":discussionId"), 10, 64)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "Invalid DiscussionId")
	}
	commentId, err := strconv.ParseInt(ctx.Params(":commentId"), 10, 64)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "Invalid CommentId")
	}
	action := ctx.Params(":action")

	req := model.DiscussionReactionRequest{
		Type:         form.Content,
		DiscussionId: discussionId,
		CommentId:    commentId,
		UserId:       ctx.Doer.ID,
	}

	switch action {
	case "react":
		// XXX: backend returns newly created reaction id, but that is useless..
		_, err := discussion_client.GiveReaction(req)
		if err != nil {
			ctx.ServerError("Failed to Give Reaction", err)
			return
		}
		log.Info("react on discussion %v's comment %v with content %v", discussionId, commentId, form.Content)
	case "unreact":
		err := discussion_client.RemoveReaction(req)
		if err != nil {
			ctx.ServerError("Failed to Remove Reaction", err)
			return
		}
		log.Info("unreact on discussion %v's comment %v with content %v", discussionId, commentId, form.Content)
	}

	// FIXME: I know this job is clumsy, but because of current backend implementation. without rewriting backend code this is the lesser of two evil..
	reactions, err := discussion_client.GetDiscussionCommentReaction(commentId)
	if err != nil {
		ctx.ServerError("Failed to Get Discussion Content", err)
		return
	}

	// i can't ensure null safety ;0
	html, err := ctx.RenderToHTML(tplDiscussionReactions, map[string]any{
		"ActionURL": fmt.Sprintf("%s/discussions/%d/comment/%d/reactions", ctx.Repo.RepoLink, discussionId, commentId),
		"Reactions": reactions.GroupByType(),
	})
	if err != nil {
		ctx.ServerError("Failed to Render Reactions..", err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{
		"html": html,
	})
}

func discussionRoleDescriptor(ctx stdCtx.Context, repo *repo_model.Repository, poster *user_model.User, discussionResponse *model.DiscussionResponse) (issues_model.RoleDescriptor, error) {
	roleDescriptor := issues_model.RoleDescriptor{}
	roleDescriptor.IsPoster = discussionResponse.IsPoster(poster.ID)

	perm, err := access_model.GetUserRepoPermission(ctx, repo, poster)
	if err != nil {
		return roleDescriptor, err
	}

	// set owner
	if perm.IsOwner() {
		roleDescriptor.RoleInRepo = issues_model.RoleRepoOwner
	}
	// set org member
	if repo.Owner.IsOrganization() {
		isMember, err := organization.IsOrganizationMember(ctx, repo.OwnerID, poster.ID)
		if err != nil {
			return roleDescriptor, err
		}
		if isMember {
			roleDescriptor.RoleInRepo = issues_model.RoleRepoMember
		}
	}
	// set collaborator
	isCollaborator, err := repo_model.IsCollaborator(ctx, repo.ID, poster.ID)
	if err != nil {
		return roleDescriptor, err
	}
	if isCollaborator {
		roleDescriptor.RoleInRepo = issues_model.RoleRepoCollaborator
	}
	// set contributor
	hasMergedPr, err := issues_model.HasMergedPullRequestInRepo(ctx, repo.ID, poster.ID)
	if err != nil {
		return roleDescriptor, err
	}
	if hasMergedPr {
		roleDescriptor.RoleInRepo = issues_model.RoleRepoContributor
	}
	return roleDescriptor, nil
}

func CreateAiSampleCodeForDiscussion(ctx *context.Context) {

	form := web.GetForm(ctx).(*structs.CreateAiSampleCodesForm)

	aiSampleCode, err := discussion.CreateAiSampleCode(
		ctx,
		&discussion.CreateDiscussionAiCommentOpt{
			Type:         form.Type,
			DiscussionId: form.DiscussionId,
			StartLine:    form.StartLine,
			EndLine:      form.EndLine,
			CodeId:       form.CodeId,
			Content:      &form.SampleCodeContent,
			GenearaterId: ctx.Doer.ID,
		})

	if err != nil {
		ctx.ServerError("ai 코드 생성 실패: %v", err)
		return
	}

	comments := make([]*discussion_service.DiscussionComment, 0, 1)
	// TODO converAiSampleCodeToDiscussionCommentFormat

	newComment, err := discussion_service.ConvertAiSampleCodeToDiscussionComment(ctx, aiSampleCode)
	if err != nil {
		ctx.ServerError(" ai 샘플코드를 디스커션 코멘트로 전환하는 과정에서 문제가 발생", err)
		return
	}

	newComment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
		Ctx: ctx,
		Links: markup.Links{
			Base: ctx.Repo.RepoLink,
		},
	}, newComment.Content)

	if err != nil {
		ctx.ServerError("markdown 렌더링 실패 : %v", err)
	}

	comments = append(comments, newComment)
	ctx.Data["DiscussionId"] = newComment.DiscussionId
	ctx.Data["comments"] = comments
	ctx.HTML(http.StatusOK, tplDiscussionFileComment)
}

func UpdateDiscussionAssignee(ctx *context.Context) {
	assigneeId := ctx.FormInt64("id")
	discussionId := ctx.FormInt64("issue_ids")
	action := ctx.FormString("action")

	switch action {
	case "clear":
		err := discussion_client.ClearDiscussionAssignee(discussionId)
		if err != nil {
			ctx.ServerError("error on discussion response: err = %v", err)
		}
	default:
		req := &model.UpdateAssigneeRequest{
			DiscussionId: discussionId,
			AssigneeId:   assigneeId,
		}
		err := discussion_client.UpdateDiscussionAssignee(req)
		if err != nil {
			ctx.ServerError("error on discussion response: err = %v", err)
		}
	}
}

func DiscussionPinOrUnpin(ctx *context.Context) {
	discussionId := ctx.ParamsInt64("discussionId")
	discussion_client.ConvertDiscussionPinStatus(discussionId)

	trimmedLink := strings.TrimSuffix(ctx.Link, "/pin")
	ctx.JSONRedirect(trimmedLink)
}

func DiscussionUnpin(ctx *context.Context) {
	discussionId := ctx.ParamsInt64("discussionId")
	discussion_client.UnpinDiscussion(discussionId)
}

func DiscussionMovePin(ctx *context.Context) {
	form := web.GetForm(ctx).(*model.MoveDiscussionPinRequest)
	discussion_client.MoveDiscussionPin(form)
}
