package structs

type CreateAiPullCommentForm struct {
	Branch       string            `json:"branch"`
	FileContents *[]PathContentMap `json:"file_contents"`
	RepoID       string            `json:"repo_id"`
	PullID       string            `json:"pull_id"`
}

type PathContentMap struct {
	TreePath string `json:"file_path"`
	Content  string `json:"code"`
}

type CreateSampleAiCommentForm struct {
	TargetCommentId string `json:"target_comment_id"`
	CommentContent  string `json:"comment"`
	CodeContent     string `json:"code"`
}
