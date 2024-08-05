package discussion

import (
	"context"

	discussion_client "code.gitea.io/gitea/client/discussion"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
)

func NewDiscussion(ctx context.Context, repo *repo_model.Repository, req *discussion_client.PostDiscussionRequest) (int, error) {
	// TODO: check poster permissions
	if user_model.IsUserBlockedBy(ctx, req.Poster, repo.OwnerID) {
		return -1, user_model.ErrBlockedUser
	}
	result, err := discussion_client.PostDiscussion(req)
	// TODO: send notification
	return result, err
}
