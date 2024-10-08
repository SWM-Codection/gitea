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

type GenerateAiSampleCodesForm struct {
	TargetCommentId string `json:"target_comment_id"`
	Type			string `json:"type"`
}

type CreateAiSampleCodesForm struct {
<<<<<<< HEAD
=======
	OriginData		  string `json:"origin_data"`
>>>>>>> 75358a09f8 (main 최신화 (#113))
	TargetCommentId   string `json:"target_comment_id"`
	SampleCodeContent string `json:"sample_code_content"`
	Type			  string `json:"type"`
}

<<<<<<< HEAD
=======

>>>>>>> 75358a09f8 (main 최신화 (#113))
type DeleteSampleCodesForm struct {
	TargetCommentId string `json:"target_comment_id"`
}

type AiSampleCodeRequest struct {
	CommentID string `json:"comment_id"`
}

type AiSampleCodeContent struct {
	ID      string  `json:"id"`
	Content *string `json:"content"`
}

type AiSampleCodeResponse struct {
	CommentID          string                 `json:"comment_id"`
	SampleCodeContent *AiSampleCodeContent `json:"content"`
}
