package model

import "code.gitea.io/gitea/modules/timeutil"

type DiscussionCommentResponse struct {
	Id           int64                 `json:"id"`
	FilePath     string                `json:"filePath"`
	GroupId      int64                 `json:"groupId"`
	DiscussionId int64                 `json:"discussionId"`
	PosterId     int64                 `json:"posterId"`
	CodeId       int64                 `json:"codeId"`
	Scope        string                `json:"scope"`
	StartLine    int64                 `json:"startLine"`
	EndLine      int64                 `json:"endLine"`
	Content      string                `json:"content"`
	CreatedUnix  timeutil.TimeStamp    `json:"createdUnix"`
	Reactions    []*DiscussionReaction `json:"reactions"`
}
