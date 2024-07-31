package discussions

import (
	"context"

	discussions_model "code.gitea.io/gitea/models/discussions"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"
)

func NewDiscussion(ctx context.Context, repo *repo_model.Repository, discussion *discussions_model.Discussion) error {
	log.Info("LoadPoster calling...")
	if err := discussion.LoadPoster(ctx); err != nil {
		return err
	}
	log.Info("discussion server info: %v", setting.DiscussionServer)

	log.Info("Discussion : %v", discussion)
	if user_model.IsUserBlockedBy(ctx, discussion.Poster, repo.OwnerID) {
		return user_model.ErrBlockedUser
	}

	log.Info("user is not blocked ok!")
	// TODO: noitfy new discussion

	return nil
}
