package discussion

import (
	"strconv"

	"code.gitea.io/gitea/services/context"

	discussion_client "code.gitea.io/gitea/client/discussion"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
)

func NewDiscussion(ctx *context.Context, repo *repo_model.Repository, req *discussion_client.PostDiscussionRequest) (int, error) {
	// TODO: check poster permissions
	if user_model.IsUserBlockedBy(ctx, req.Poster, repo.OwnerID) {
		return -1, user_model.ErrBlockedUser
	}
	result, err := discussion_client.PostDiscussion(req)
	// TODO: send notification
	return result, err
}

func GetDiscussionList(ctx *context.Context) (*discussion_client.DiscussionListResponse, error) {
	repo := ctx.Repo.Repository
	repoId := repo.ID
	isClosed := ctx.FormString("state") == "closed"
	page, err := strconv.Atoi(ctx.FormString("page"))
	if err != nil {
		page = 1
	}
	page-- // gitea uses 1-based paginiation methodology, but discussion backend uses 0-based pagination methodology

	discussionListResponse, err := discussion_client.GetDiscussionList(repoId, isClosed, page)
	if err != nil {
		return nil, err
	}
	// post process discussions
	for _, d := range discussionListResponse.Discussions {
		d.LoadRepo(ctx)
		d.LoadPoster(ctx)
	}
	return discussionListResponse, nil
}
