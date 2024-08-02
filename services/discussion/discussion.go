package discussion

import (
	"context"

	discussion_client "code.gitea.io/gitea/client/discussion"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
)

func NewDiscussion(ctx context.Context, repo *repo_model.Repository, req *discussion_client.PostDiscussionRequest) error {
	// TODO: check poster

	if user_model.IsUserBlockedBy(ctx, req.Poster, repo.OwnerID) {
		return user_model.ErrBlockedUser
	}
	// TODO: notify new discussion
	return nil
}
