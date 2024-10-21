package model

type DiscussionCode struct {
	Id        int64  `json:"id"`
	FilePath  string `json:"filePath"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}
