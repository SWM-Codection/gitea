package model

import user_model "code.gitea.io/gitea/models/user"

type PostDiscussionRequest struct {
	RepoId     int64            `json:"repoId"`
	Poster     *user_model.User `json:"-"`
	PosterId   int64            `json:"posterId"`
	Name       string           `json:"name"`
	Content    string           `json:"content"`
	BranchName string           `json:"branchName"`
	Codes      []DiscussionCode `json:"codes"`
}
