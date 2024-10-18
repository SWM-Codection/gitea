package model

type ExtractedLine struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
}
