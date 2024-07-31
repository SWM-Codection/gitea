package discussions

import (
	"context"
	"fmt"

	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
)

type Discussion struct {
	ID       int64
	RepoID   int64
	Repo     *repo_model.Repository `xorm:"-"`
	Title    string
	PosterID int64
	Poster   *user_model.User `xorm:"-"`
	Content  string
}

func (discussion *Discussion) LoadPoster(ctx context.Context) (err error) {
	if discussion.Poster == nil && discussion.PosterID != 0 {
		discussion.Poster, err = user_model.GetPossibleUserByID(ctx, discussion.PosterID)
		if err != nil {
			discussion.PosterID = user_model.GhostUserID
			discussion.Poster = user_model.NewGhostUser()
			if !user_model.IsErrUserNotExist(err) {
				return fmt.Errorf("getUserByID.(poster) [%d]: %w", discussion.PosterID, err)
			}
			return nil
		}
	}
	return err
}
