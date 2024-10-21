package model

import (
	"time"

	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/timeutil"
)

type Discussion struct {
	Id           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Content      string                 `json:"content"`
	RepoId       int64                  `json:"repoId"`
	PosterId     int64                  `json:"posterId"`
	CommitHash   string                 `json:"commitHash"`
	Index        int64                  `json:"index"`
	IsClosed     bool                   `json:"isClosed"`
	CreatedUnix  timeutil.TimeStamp     `json:"createdUnix"` // required, but didn't exist before
	ClosedUnix   timeutil.TimeStamp     `json:"closedUnix"`  // required, but didn't exist before
	DeadlineUnix timeutil.TimeStamp     `json:"deadlineUnix"`
	NumComments  int                    `json:"-"` // it can be computed
	Repo         *repo_model.Repository `json:"-"` // it can be computed via d.LoadRepo()
	Poster       *user_model.User       `json:"-"` // it can be computed via d.LoadPoster()
}

/**
 * discussion methods
 * the `discussion` struct could be moved to a separate file later
 */
func (d *Discussion) IsExpired() bool {
	return d.DeadlineUnix < timeutil.TimeStamp(time.Now().Unix())
}

func (d *Discussion) GetLastEventTimestamp() timeutil.TimeStamp {
	if d.IsClosed {
		return d.ClosedUnix
	}
	return d.CreatedUnix
}

func (d *Discussion) GetLastEventLabel() string {
	if d.IsClosed {
		return "repo.discussion.closed_by"
	}
	return "repo.discussion.opened_by"
}

func (d *Discussion) GetLastEventLabelFake() string {
	if d.IsClosed {
		return "repo.discussion.closed_by_fake"
	}
	return "repo.discussion.opened_by_fake"
}

// func (d *Discussion) LoadRepo(ctx *context.Context) error {

// 	repo, err := repo.GetRepositoryByID(*ctx, d.RepoId)
// 	if err != nil {
// 		log.Printf("Error getting repo by id: %v", d.RepoId)
// 		return err
// 	}
// 	d.Repo = repo
// 	return nil
// }

// func (d *Discussion) LoadPoster(ctx *context.Context) error {
// 	poster, err := user_model.GetUserByID(*ctx, d.PosterId)
// 	if err != nil {
// 		log.Printf("Error getting repo by id: %v", d.PosterId)
// 		return err
// 	}
// 	d.Poster = poster
// 	return nil
// }
