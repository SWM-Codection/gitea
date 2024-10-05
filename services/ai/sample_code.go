package ai

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"code.gitea.io/gitea/client/discussion"
	discussion_model "code.gitea.io/gitea/models/discussion"
	issues_model "code.gitea.io/gitea/models/issues"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/setting"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"github.com/spf13/cast"
)

type SampleCodeService interface {
	GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64, sampleType string) (*api.AiSampleCodeResponse, error)
	GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*GenerateSampleCodeResponse, error)
	CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm) (*discussion_model.AiSampleCode, error)
	DeleteAiSampleCode(ctx *context.Context, id int64) error
}

var DEFAULT_CAPACITY int64 = 10

type SampleCodeServiceImpl struct{}

var _ SampleCodeService = &SampleCodeServiceImpl{}

func (is *SampleCodeServiceImpl) GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64, sampleType string) (*api.AiSampleCodeResponse, error) {

	response, err := AiSampleCodeDbAdapter.GetAiSampleCodesByCommentID(ctx, commentID, sampleType)

	if err != nil {
		return nil, err
	}

	return response, nil

}

func (is *SampleCodeServiceImpl) CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm) (*discussion_model.AiSampleCode, error) {
	// TODOC Discussion Comment 무결성 검사

	aiSampleCode, err := AiSampleCodeDbAdapter.InsertAiSampleCode(ctx, &discussion_model.CreateDiscussionAiCommentOpt{
		TargetCommentId: cast.ToInt64(form.TargetCommentId),
		GenearaterId:    ctx.Doer.ID,
		Type:            form.Type,
		Content:         &form.SampleCodeContent,
	})

	if err != nil {
		return nil, err
	}

	return aiSampleCode, nil

}

func GetFileContentFromCommit(ctx *context.Context, repoPath, commitHash, filePath string) (string, error) {
	// 레포지토리 열기
	repo, err := git.OpenRepository(ctx, repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}
	defer repo.Close()

	// 커밋 객체 가져오기
	commit, err := repo.GetCommit(commitHash)
	if err != nil {
		return "", fmt.Errorf("failed to get commit: %w", err)
	}

	// 파일 경로로 Blob 객체 가져오기
	entry, err := commit.GetTreeEntryByPath(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get tree entry: %w", err)
	}
	blob := entry.Blob()

	// Blob의 내용 가져오기
	content, err := blob.GetBlobContent(setting.UI.MaxDisplayFileSize) // limit은 가져올 내용의 바이트 수
	if err != nil {
		return "", fmt.Errorf("failed to get blob content: %w", err)
	}

	return content, nil
}

func getContentForFull(ctx *context.Context, targetCommentId int64, form *api.GenerateAiSampleCodesForm) (*string, *string, error) {

	comment, err := issues_model.GetCommentByID(ctx, targetCommentId)
	if err != nil {
		return nil, nil, fmt.Errorf("Comment not found: %v", err)
	}

	issue, err := issues_model.GetIssueByID(ctx, comment.IssueID)
	if err != nil {
		return nil, nil, fmt.Errorf("Issue not found: %v", err)
	}

	repo, err := repo_model.GetRepositoryByID(ctx, issue.RepoID)
	if err != nil {
		return nil, nil, fmt.Errorf("Repository not found: %v", err)
	}

	codeContent, err := GetFileContentFromCommit(ctx, repo.RepoPath(), comment.CommitSHA, comment.TreePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get file content from commit: %v", err)
	}

	return &codeContent, &comment.Content, nil
}

func getContentForDiscussion(ctx *context.Context, targetCommentId int64, form *api.GenerateAiSampleCodesForm) (*string, *string, error) {
	comment, err := discussion.GetDiscussionComment(targetCommentId)

	if err != nil {
		return nil, nil, fmt.Errorf("Comment found error: %v", err)
	}

	discussion, err := discussion.GetDiscussion(comment.DiscussionId)
	if err != nil {
		return nil, nil, fmt.Errorf("Discussion found error: %v", err)
	}
	repo, err := repo_model.GetRepositoryByID(ctx, discussion.RepoId)
	if err != nil {
		return nil, nil, fmt.Errorf("repo found error: %v", err)
	}

	fileContent, err := GetFileContentFromCommit(ctx, repo.RepoPath(), discussion.CommitHash, comment.FilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get file content from commit: %v", err)
	}

	lines := strings.Split(fileContent, "\n")

	if comment.StartLine < 1 || comment.StartLine > comment.EndLine {
		return nil, nil, fmt.Errorf("Invalid start or end line")
	}

	codeContent := strings.Join(lines[comment.StartLine-1:comment.EndLine], "\n")

	return &codeContent, &comment.Content, nil
}

func (is *SampleCodeServiceImpl) GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*GenerateSampleCodeResponse, error) {
	targetCommentId, err := strconv.ParseInt(form.TargetCommentId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid TargetCommentId: %v", err)
	}

	response, err := AiSampleCodeDbAdapter.GetAiSampleCodesByCommentID(ctx, targetCommentId, form.Type)

	if err != nil || response.SampleCodeContent != nil {
		return nil, fmt.Errorf("already Ai comment: %v", err)
	}

	var codeContent, commentContent *string

	if form.Type == "pull" {
		codeContent, commentContent, err = getContentForFull(ctx, targetCommentId, form)

	} else if form.Type == "discussion" {

		codeContent, commentContent, err = getContentForDiscussion(ctx, targetCommentId, form)

	}

	if err != nil {
		return nil, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(AI_SAMPLE_CODE_UNIT)

	resultQueue := make(chan *GenerateSampleCodeResponse, AI_SAMPLE_CODE_UNIT)

	for i := 0; i < AI_SAMPLE_CODE_UNIT; i++ {
		go func(codeContent, commentContent string) {
			defer wg.Done()

			response, err := AiSampleCodeRequester.RequestReviewToAI(ctx, &AiSampleCodeRequest{
				CodeContent:    codeContent,
				CommentContent: commentContent,
			})

			if err != nil {

				fmt.Errorf("request sample to ai server fail")
				resultQueue <- nil
				return
			}

			result := &GenerateSampleCodeResponse{
				SampleCode:       "",
				OriginalMarkdown: response.SampleCode,
			}

			highlightedCode, _ := markup.RenderString(&markup.RenderContext{
				Ctx:  git.DefaultContext,
				Type: "markdown",
			}, result.OriginalMarkdown)

			result.SampleCode = string(highlightedCode)

			resultQueue <- result
		}(*codeContent, *commentContent)
	}

	wg.Wait()
	close(resultQueue)

	sampleCodes := make([]*GenerateSampleCodeResponse, 0, AI_SAMPLE_CODE_UNIT)
	for result := range resultQueue {
		sampleCodes = append(sampleCodes, result)
	}

	return sampleCodes, nil
}

func (is *SampleCodeServiceImpl) DeleteAiSampleCode(ctx *context.Context, id int64) error {

	return AiSampleCodeDbAdapter.DeleteAiSampleCodeByID(ctx, id)
}
