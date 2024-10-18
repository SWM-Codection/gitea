package model

type ModifyDiscussionCommentRequest struct {
	DiscussionId        int64            `json:"discussionId"`
	DiscussionCommentId int64            `json:"discussionCommentId"`
	CodeId              *int64           `json:"codeId"`
	PosterId            int64            `json:"posterId"`
	Scope               CommentScopeEnum `json:"scope"`
	StartLine           *int32           `json:"startLine"`
	EndLine             *int32           `json:"endLine"`
	Content             string           `json:"content"`
}
