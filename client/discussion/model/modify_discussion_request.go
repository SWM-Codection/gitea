package model

type ModifyDiscussionRequest struct {
	RepoId       int64            `json:"repoId"`
	DiscussionId int64            `json:"discussionId"`
	PosterId     int64            `json:"posterId"`
	Name         string           `json:"name"`
	Content      string           `json:"content"`
	Codes        []DiscussionCode `json:"codes"`
}
