package model

type DiscussionReaction struct {
	Id           int64            `json:"id"`
	Type         ReactionTypeEnum `json:"type"`
	DiscussionId int64            `json:"discussionId"`
	CommentId    int64            `json:"commentId"`
	UserId       int64            `json:"userId"`
}
