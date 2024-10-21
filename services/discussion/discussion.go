package discussion

import (
	"fmt"
	"strconv"

	"code.gitea.io/gitea/modules/highlight"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/services/context"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/client/discussion/model"
	repo_model "code.gitea.io/gitea/models/repo"
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
		d.LoadRepo(ctx)
		d.LoadPoster(ctx)
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
		d.LoadRepo(ctx)
		d.LoadPoster(ctx)
	}
	return discussionListResponse, nil
}
