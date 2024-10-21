package model

type DeleteDiscussionCommentRequest struct {
	PosterId            int64 `json:"posterId"`
	DiscussionCommentId int64 `json:"discussionCommentId"`
}
