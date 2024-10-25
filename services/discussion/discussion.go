package discussion

import (
	"fmt"
	"html/template"
	"strconv"

	"code.gitea.io/gitea/models/discussion"

	"code.gitea.io/gitea/modules/highlight"
	"code.gitea.io/gitea/modules/timeutil"
	"code.gitea.io/gitea/modules/util"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/client/discussion/model"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/user"
	user_model "code.gitea.io/gitea/models/user"
)

func NewDiscussion(ctx *context.Context, repo *repo_model.Repository, req *model.PostDiscussionRequest) (int, error) {
	// TODO: check poster permissions
	if user_model.IsUserBlockedBy(ctx, req.Poster, repo.OwnerID) {
		return -1, user_model.ErrBlockedUser
	}
	result, err := discussion_client.PostDiscussion(req)
	// TODO: send notification
	return result, err
}

func GetDiscussionList(ctx *context.Context) (*model.DiscussionListResponse, error) {
	repo := ctx.Repo.Repository
	repoId := repo.ID
	isClosed := ctx.FormString("state") == "closed"
	page, err := strconv.Atoi(ctx.FormString("page"))
	if err != nil {
		page = 1
	}
	page-- // gitea uses 1-based paginiation methodology, but discussion backend uses 0-based pagination methodology
	sort := ctx.FormString("sort")

	log.Info("calling get discussion list with repoId: %v, isClosed: %v, page: %v", repoId, isClosed, page)
	discussionListResponse, err := discussion_client.GetDiscussionList(repoId, isClosed, page, sort)
	if err != nil {
		log.Info("discusisonClient.getdiscussionList failed1")
		return nil, err
	}
	// post process discussions
	for _, d := range discussionListResponse.Discussions {
		d.Repo = ctx.Repo.Repository
		poster, _ := user.GetUserByID(ctx, d.PosterId)
		d.Poster = poster
	}
	return discussionListResponse, nil
}

func GetDiscussionContent(discussionID int64) (*model.DiscussionContentResponse, error) {
	contents, err := discussion_client.GetDiscussionContents(discussionID)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func GetDiscussionContentWithHighlights(discussionID int64) (*model.DiscussionContentResponse, error) {
	contents, err := discussion_client.GetDiscussionContents(discussionID)
	if err != nil {
		return nil, err
	}

	err = highlightContents(contents)

	if err != nil {
		return nil, err
	}

	return contents, nil

}

func highlightContents(contents *model.DiscussionContentResponse) error {
	for i := range contents.Contents {
		if err := highlightContent(&contents.Contents[i]); err != nil {
			return err
		}
	}
	return nil
}

func highlightContent(content *model.FileContent) error {
	for i, block := range content.CodeBlocks {
		for j, line := range block.Lines {
			html, _, err := highlight.File(content.FilePath, "", []byte(line.Content))
			if err != nil {
				return fmt.Errorf("failed to highlight line %d in block %d: %w", j, i, err)
			}
			if len(html) == 0 {
				continue
			}

			content.CodeBlocks[i].Lines[j].Content = string(html[0])
		}
	}
	return nil
}

func DeleteDiscussionComment(ctx *context.Context, discussionId int64, posterId int64) error {

	if err := discussion_client.DeleteDiscussionComment(discussionId, posterId); err != nil {
		return err
	}

	return nil
}

func GetDiscussionCommentsByCodeId(ctx *context.Context, codeId int64) ([]*DiscussionComment, error) {

	commentsResp, err := discussion_client.GetDiscussionCommentsByCodeId(codeId)

	sampleCodes, err := discussion.GetAiSampleCodesByCodeId(ctx, codeId, "discussion")

	comments := make([]*DiscussionComment, 0, len(sampleCodes))

	for _, sampleCode := range sampleCodes {

		aiComment, err := ConvertAiSampleCodeToDiscussionComment(ctx, sampleCode)

		if err != nil {
			ctx.ServerError("failed to convert ai sample code to discussion comment: %v", err)
		}
		comments = append(comments, aiComment)
	}

	for _, comment := range commentsResp {
		poster, err := user_model.GetUserByID(ctx, comment.PosterId)

		if err != nil {

		}
		discussionComment := &DiscussionComment{
			ID:           comment.Id,
			StartLine:    comment.StartLine,
			PosterId:     comment.PosterId,
			DiscussionId: comment.DiscussionId,
			CodeId:       comment.CodeId,
			GroupId:      comment.GroupId,
			EndLine:      comment.EndLine,
			CreatedUnix:  comment.CreatedUnix,
			Reactions:    comment.Reactions,
			Poster:       poster,
			Content:      comment.Content,
		}
		comments = append(comments, discussionComment)
	}

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func ModifyDiscussionComment(ctx *context.Context, form *forms.ModifyDiscussionCommentForm) error {

	commentScope := util.Iif(
		form.CodeId != nil && form.StartLine != nil && form.EndLine != nil,
		model.CommentScopeLocal,
		model.CommentScopeGlobal,
	)

	posterId := ctx.Doer.ID

	discussionId := ctx.ParamsInt64(":discussionId")

	request := &model.ModifyDiscussionCommentRequest{
		DiscussionId:        discussionId,
		DiscussionCommentId: form.DiscussionCommentId,
		CodeId:              form.CodeId,
		PosterId:            posterId,
		Scope:               commentScope,
		StartLine:           form.StartLine,
		EndLine:             form.EndLine,
		Content:             form.Content,
	}

	err := discussion_client.ModifyDiscussionComment(request)

	if err != nil {
		return err
	}

	return nil
}

func ConvertAiSampleCodeToDiscussionComment(ctx *context.Context, sampleCode *discussion.AiSampleCode) (*DiscussionComment, error) {

	aiPoster, err := user_model.GetPossibleUserByID(ctx, -3)

	if err != nil {
		return nil, err
	}

	newComment := &DiscussionComment{
		ID:           -sampleCode.Id,
		StartLine:    sampleCode.StartLine,
		DiscussionId: sampleCode.DiscussionId,
		GroupId:      -sampleCode.Id,
		EndLine:      sampleCode.EndLine,
		CodeId:       sampleCode.CodeId,
		CreatedUnix:  sampleCode.CreatedUnix,
		Reactions:    nil, // TODO: 뱃지 형식으로 변경하기
		Poster:       aiPoster,
		Content:      sampleCode.Content,
		PosterId:     sampleCode.GenearaterId,
	}

	return newComment, err
}

type DiscussionComment struct {
	ID              int64
	DiscussionId    int64
	PosterId        int64
	Poster          *user_model.User
	GroupId         int64
	Content         string
	StartLine       int64
	CodeId          int64
	EndLine         int64
	Reactions       model.ReactionList
	RenderedContent template.HTML
	CreatedUnix     timeutil.TimeStamp
}

func (c *DiscussionComment) HashTag() string {
	return fmt.Sprintf("discussioncomment-%d", c.ID)
}

func (c *DiscussionComment) IsAiSampleCode() bool {
	return c.Poster.ID == -3
}

func GetPinnedDiscussionList(ctx *context.Context) (*model.DiscussionListResponse, error) {
	repo := ctx.Repo.Repository
	repoId := repo.ID
	discussionListResponse, err := discussion_client.GetPinnedDiscussions(repoId)
	if err != nil {
		log.Error("discussionClient.getPinnedDiscussions failed")
		return nil, err
	}
	// post process discussions
	for _, d := range discussionListResponse.Discussions {
		d.Repo = repo
		poster, _ := user_model.GetUserByID(ctx, d.PosterId)
		d.Poster = poster
	}
	return discussionListResponse, nil
}
