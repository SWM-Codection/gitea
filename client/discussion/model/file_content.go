package model

type FileContent struct {
	FilePath   string      `json:"filePath"`
	CodeBlocks []CodeBlock `json:"codeBlocks"`
}
