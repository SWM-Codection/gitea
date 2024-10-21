package model

type DiscussionListResponse struct {
	TotalCount  int64         `json:"totalCount"`
	Discussions []*Discussion `json:"discussions"`
}
