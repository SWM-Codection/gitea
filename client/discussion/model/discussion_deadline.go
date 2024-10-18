package model

import "time"

type DiscussionDeadline struct {
	Deadline *time.Time `json:"due_date"`
}
