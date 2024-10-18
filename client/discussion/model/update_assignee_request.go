package model

type UpdateAssigneeRequest struct {
	DiscussionId int64 `json:"discussionId"`
	AssigneeId   int64 `json:"assigneeId"`
}
