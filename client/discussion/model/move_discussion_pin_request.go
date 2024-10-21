package model

type MoveDiscussionPinRequest struct {
	DiscussionId int64 `json:"id"`
	Position     int64 `json:"position"`
}
