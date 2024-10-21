package model

import "code.gitea.io/gitea/modules/timeutil"

type DiscussionErrorResponse struct {
	TimeStamp timeutil.TimeStamp `json:"timestamp"`
	Status    int                `json:"status"`
	Error     string             `json:"error"`
	Message   string
}
