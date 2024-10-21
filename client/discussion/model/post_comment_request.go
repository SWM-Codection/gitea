package model

type PostCommentRequest struct {
	DiscussionId int64            `json:"discussionId"`
	CodeId       *int64           `json:"codeId"`
	GroupId      *int64           `json:"groupId"`
	PosterId     int64            `json:"posterId"`
	Scope        CommentScopeEnum `json:"scope"`
	StartLine    *int32           `json:"startLine"`
	EndLine      *int32           `json:"endLine"`
	Content      string           `json:"content"`
}
