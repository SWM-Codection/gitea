package model

import (
	"code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/timeutil"
)

type DiscussionResponse struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	Content     string             `json:"content"`
	RepoId      int64              `json:"repoId"`
	PosterId    int64              `json:"posterId"`
	CommitHash  string             `json:"commitHash"`
	IsClosed    bool               `json:"isClosed"`
	Deadline    timeutil.TimeStamp `json:"deadline"`
	Assignees   []int64            `json:"assignees"`
	CreatedUnix timeutil.TimeStamp `json:"createdUnix"`
	UpdatedUnix timeutil.TimeStamp `json:"updatedUnix"`
	Index       int64              `json:"index"`
	Poster      *user.User         `json:"-"`
}

func (dr DiscussionResponse) IsPoster(id int64) bool {
	return dr.PosterId == id
}
