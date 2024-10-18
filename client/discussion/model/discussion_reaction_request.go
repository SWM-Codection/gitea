package model

type DiscussionReactionRequest struct {
	Type         ReactionTypeEnum `json:"type"`
	DiscussionId int64            `json:"discussionId"`
	CommentId    int64            `json:"commentId"`
	UserId       int64            `json:"userId"`
}
