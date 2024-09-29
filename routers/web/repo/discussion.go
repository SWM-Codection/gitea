package repo

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	}

	req := &discussion_client.PostCommentRequest{
		DiscussionId: discussionId,
		Scope:        commentScope,
		PosterId:     ctx.Doer.ID,
		Content:      form.Content,
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

type ReactionList []*discussion_client.DiscussionReaction

func (list ReactionList) GroupByType() map[string]ReactionList {
	reactions := make(map[string]ReactionList)
	for _, reaction := range list {
		reactions[reaction.Type] = append(reactions[reaction.Type], reaction)
	}
	return reactions
}

type DiscussionComment struct {
	ID              int64
	Poster          *user_model.User
	Content         string
	StartLine       int64
	EndLine         int64
	Reactions       ReactionList
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
		ID:          comment.Id,
		StartLine:   comment.StartLine,
		EndLine:     comment.EndLine,
		CreatedUnix: comment.CreatedUnix,
		Reactions:   comment.Reactions,
		Poster:      poster,
		Content:     comment.Content,
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

	ctx.HTML(http.StatusOK, tplDiscussionFileComments)

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
	println("here")
	println(discussionId)

	var deadlineUnix int64 // Unix 타임스탬프를 저장할 int64 변수
	var deadline time.Time
	if form.Deadline != nil && !form.Deadline.IsZero() {
    	deadline = time.Date(form.Deadline.Year(), form.Deadline.Month(), form.Deadline.Day(),
        	23, 59, 59, 0, time.Local)
    	deadlineUnix = deadline.Unix() // Unix 타임스탬프를 int64로 저장
	}
	println(deadlineUnix)

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
