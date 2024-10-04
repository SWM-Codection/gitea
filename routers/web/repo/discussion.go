package repo

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"code.gitea.io/gitea/models/discussion"

	stdCtx "context"
	api "code.gitea.io/gitea/modules/structs"

	discussion_client "code.gitea.io/gitea/client/discussion"
	issues_model "code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/models/organization"
	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/markup/markdown"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/modules/timeutil"
	"code.gitea.io/gitea/modules/util"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	discussion_service "code.gitea.io/gitea/services/discussion"
	"code.gitea.io/gitea/services/forms"
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

	log.Info("discussions : %v", listResp.Discussions)
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
		println(len(assignees))
	}
	repo := ctx.Repo.Repository
	assigneeUsers, err := repo_model.GetRepoAssignees(ctx, repo)
	if err != nil {
		ctx.ServerError("GetRepoAssignees", err)
		return
	}
	println(len(assignees))

	participants[0] = poster
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
	ctx.PageData["RepoLink"] = ctx.Repo.Repository.RepoPathLink()
	ctx.PageData["DiscussionId"] = discussionId
	ctx.Data["DiscussionTab"] = "files"

	ctx.HTML(http.StatusOK, tplDiscussionView)
}

func NewDiscussionCommentPost(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.CreateDiscussionCommentForm)

	discussionId := ctx.ParamsInt64(":discussionId")
	commentScope := util.Iif(
		form.CodeId != nil && form.StartLine != nil && form.EndLine != nil,
		discussion_client.CommentScopeLocal,
		discussion_client.CommentScopeGlobal,
	)
	// if request malformed, then make it valid
	if commentScope == discussion_client.CommentScopeGlobal {
		form.CodeId = nil
		form.StartLine = nil
		form.EndLine = nil
		form.GroupId = nil
	}

	req := &discussion_client.PostCommentRequest{
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

type DiscussionComment struct {
	ID              int64
	DiscussionId    int64
	Poster          *user_model.User
	GroupId         int64
	Content         string
	StartLine       int64
	CodeId          int64
	EndLine         int64
	Reactions       discussion_client.ReactionList
	RenderedContent template.HTML
	CreatedUnix     timeutil.TimeStamp
}

func (c *DiscussionComment) HashTag() string {
	return fmt.Sprintf("discussioncomment-%d", c.ID)
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
	comments := make([]*DiscussionComment, 0, 1)

	newComment := &DiscussionComment{
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
	commentsResp []*discussion_client.DiscussionCommentResponse) (map[int64][]*DiscussionComment, error) {

	comments := make(map[int64][]*DiscussionComment)

	for _, comment := range commentsResp {
		poster, err := user_model.GetUserByID(ctx, comment.PosterId)

		if err != nil {
			return nil, err
		}

		discussionComment := DiscussionComment{
			ID:           comment.Id,
			StartLine:    comment.StartLine,
			DiscussionId: comment.DiscussionId,
			CodeId:       comment.CodeId,
			GroupId:      comment.GroupId,
			EndLine:      comment.EndLine,
			CreatedUnix:  comment.CreatedUnix,
			Reactions:    comment.Reactions,
			Poster:       poster,
			Content:      comment.Content,
		}

		discussionComment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: ctx.Repo.RepoLink,
			},
		}, discussionComment.Content)

		if err != nil {
			return nil, err
		}

		comments[discussionComment.GroupId] = append(comments[discussionComment.GroupId], &discussionComment)

		sampleCode, err := discussion.GetAiSampleCodeByCommentID(ctx, comment.Id, "discussion")
		if err != nil {
			return nil, err
		}
		if sampleCode == nil {
			continue
		}

		aiComment, err := convertAiSampleCodeToDiscussionComment(ctx, sampleCode, comment)

		aiComment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: ctx.Repo.RepoLink,
			},
		}, aiComment.Content)

		comments[aiComment.GroupId] = append(comments[aiComment.GroupId], aiComment)
	}

	return comments, nil
}

func renderDiscussionFileComments(ctx *context.Context, commentGroups map[int64][]*DiscussionComment) ([]*DiscussionFileCommentsResponse, error) {

	resp := make([]*DiscussionFileCommentsResponse, 0, len(commentGroups))

	for _, groupComments := range commentGroups {

		sort.Slice(groupComments, func(i, j int) bool {
			return groupComments[i].CreatedUnix > groupComments[j].CreatedUnix
		})

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

	commentsResp, err := discussion_client.GetDiscussionCommentsByCodeId(codeId)

	if err != nil {
		ctx.JSONError(err)
	}

	commentGroups, err := groupDiscussionFileCommentsByGroupId(ctx, commentsResp)

	if err != nil {
		ctx.JSONError(err)
	}

	resp, err := renderDiscussionFileComments(ctx, commentGroups)

	if err != nil {
		ctx.JSONError(err)
	}

	ctx.JSON(http.StatusOK, resp)

}

// RenderNewCodeCommentForm will render the form for creating a new review comment
func RenderNewDiscussionFileCommentForm(ctx *context.Context) {

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

	ctx.JSON(http.StatusCreated, discussion_client.DiscussionDeadline{Deadline: &deadline})
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
	ctx.JSONOK()
}

func DeleteDiscussionFileComment(ctx *context.Context) {
	posterId := ctx.Doer.ID

	discussionCommentId := ctx.ParamsInt64(":discussionId")

	err := discussion_service.DeleteDiscussionComment(ctx, discussionCommentId, posterId)

	if err != nil {
		log.Error(err.Error())
		ctx.JSONError(err.Error())
	}

	ctx.JSONOK()

}

func ModifyDiscussionFileComment(ctx *context.Context) {

	form := web.GetForm(ctx).(*forms.ModifyDiscussionCommentForm)

	commentScope := util.Iif(
		form.CodeId != nil && form.StartLine != nil && form.EndLine != nil,
		discussion_client.CommentScopeLocal,
		discussion_client.CommentScopeGlobal,
	)

	posterId := ctx.Doer.ID

	discussionId := ctx.ParamsInt64(":discussionId")

	request := &discussion_client.ModifyDiscussionCommentRequest{
		DiscussionId:        discussionId,             // discussionId 변수는 이미 전달된 것으로 가정
		DiscussionCommentId: form.DiscussionCommentId, // form에서 넘어온 DiscussionCommentId
		CodeId:              form.CodeId,              // form에서 넘어온 CodeId (포인터)
		PosterId:            posterId,                 // posterId는 전달된 값 (별도의 변수로 가정)
		Scope:               commentScope,             // scope는 전달된 값 (별도의 변수로 가정, CommentScopeEnum 타입)
		StartLine:           form.StartLine,           // form에서 넘어온 StartLine (포인터)
		EndLine:             form.EndLine,             // form에서 넘어온 EndLine (포인터)
		Content:             form.Content,             // form에서 넘어온 Content
	}

	err := discussion_client.ModifyDiscussionComment(request)

	if err != nil {
		ctx.ServerError("modify request failed %v", err)
		return
	}

	var renderedContent template.HTML
	if request.Content != "" {
		renderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: ctx.Repo.RepoLink,
			},
		}, request.Content)
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

	req := discussion_client.DiscussionReactionRequest{
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
		}
		log.Info("react on discussion %v's comment %v with content %v", discussionId, commentId, form.Content)
	case "unreact":
		err := discussion_client.RemoveReaction(req)
		if err != nil {
			ctx.ServerError("Failed to Remove Reaction", err)
		}
		log.Info("unreact on discussion %v's comment %v with content %v", discussionId, commentId, form.Content)
	}

	// FIXME: I know this job is clumsy, but because of current backend implementation. without rewriting backend code this is the lesser of two evil..
	d, err := discussion_client.GetDiscussionContent(discussionId)
	if err != nil {
		ctx.ServerError("Failed to Get Discussion Content", err)
	}
	var reactions discussion_client.ReactionList
	for _, gc := range d.GlobalComments {
		if gc.Id == commentId {
			reactions = gc.Reactions
		}
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

func discussionRoleDescriptor(ctx stdCtx.Context, repo *repo_model.Repository, poster *user_model.User, discussionResponse *discussion_client.DiscussionResponse) (issues_model.RoleDescriptor, error) {
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

	commentId, err := strconv.ParseInt(form.TargetCommentId, 10, 64)
	if err != nil {
		ctx.ServerError("잘못된 CommentId 형식", err)
	}

	aiSampleCode, err := discussion.CreateAiSampleCode(
		ctx,
		&discussion.CreateDiscussionAiCommentOpt{
			Type:            form.Type,
			TargetCommentId: commentId,
			Content:         &form.SampleCodeContent,
			GenearaterId:    ctx.Doer.ID,
		})

	if err != nil {
		ctx.ServerError("ai 코드 생성 실패: %v", err)
		return
	}

	comments := make([]*DiscussionComment, 0, 1)
	// TODO converAiSampleCodeToDiscussionCommentFormat
	comment, err := discussion_client.GetDiscussionComment(aiSampleCode.TargetCommentId)

	newComment, err := convertAiSampleCodeToDiscussionComment(ctx, aiSampleCode, comment)
	if err != nil {
		ctx.ServerError("문제가 발생", err)
	}

	newComment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
		Ctx: ctx,
		Links: markup.Links{
			Base: ctx.Repo.RepoLink,
		},
	}, newComment.Content)

	comments = append(comments, newComment)
	ctx.Data["DiscussionId"] = newComment.DiscussionId
	ctx.Data["comments"] = comments
	ctx.HTML(http.StatusOK, tplDiscussionFileComment)
}

func convertAiSampleCodeToDiscussionComment(ctx *context.Context, sampleCode *discussion.AiSampleCode, comment *discussion_client.DiscussionCommentResponse) (*DiscussionComment, error) {

	poster, err := user_model.GetPossibleUserByID(ctx, -3)

	if err != nil {
		return nil, err
	}

	newComment := &DiscussionComment{
		ID:           -comment.Id,
		StartLine:    comment.StartLine,
		DiscussionId: comment.DiscussionId,
		GroupId:      comment.GroupId,
		EndLine:      comment.EndLine,
		CodeId:       comment.CodeId,
		CreatedUnix:  comment.CreatedUnix - 1,
		Reactions:    nil, // TODO: 뱃지 형식으로 변경하기
		Poster:       poster,
		Content:      sampleCode.Content,
	}

	return newComment, err
}


func UpdateDiscussionAssignee(ctx *context.Context)  {
	assigneeId := ctx.FormInt64("id")
	discussionId := ctx.FormInt64("issue_ids")
	action := ctx.FormString("action")
	println(assigneeId)
	println(discussionId)
	println(action)

	switch action {
	case "clear":
		err := discussion_client.ClearDiscussionAssignee(discussionId)
		if err != nil {
			ctx.ServerError("error on discussion response: err = %v", err)
		}
	default:
		req := &discussion_client.UpdateAssigneeRequest{
			DiscussionId:	discussionId,
			AssigneeId: 	assigneeId,
		}
		err := discussion_client.UpdateDiscussionAssignee(req)
		if err != nil {
			ctx.ServerError("error on discussion response: err = %v", err)
		}
	}
}