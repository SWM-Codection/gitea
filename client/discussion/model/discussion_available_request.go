package model

type DiscussionAvailableRequest struct {
	RepoId    int64 `json:"repoId"`
	Available bool  `json:"available"`
}
