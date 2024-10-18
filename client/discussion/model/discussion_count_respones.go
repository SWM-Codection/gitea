package model

type DiscussionCountResponse struct {
	OpenCount  int `json:"openCount"`
	CloseCount int `json:"closeCount"`
}
