package model

type DiscussionContentResponse struct {
	DiscussionId    int64                       `json:"discussionId"`
	Contents        []FileContent               `json:"contents"`
	GlobalComments  []DiscussionCommentResponse `json:"globalComments"`
	GlobalReactions ReactionList                `json:"discussionReaction"`
}
