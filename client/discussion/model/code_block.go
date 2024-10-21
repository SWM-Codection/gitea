package model

type CodeBlock struct {
	CodeId   int64                       `json:"codeId"`
	Lines    []ExtractedLine             `json:"lines"`
	Comments []DiscussionCommentResponse `json:"comments"`
}
