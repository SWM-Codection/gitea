package ai

import (
	"strconv"
	discussion_model "code.gitea.io/gitea/models/discussion"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"fmt"
	"github.com/spf13/cast"
	"sync"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/setting"
	issues_model "code.gitea.io/gitea/models/issues"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/modules/highlight"
)

type SampleCodeService interface {
	GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error)
	GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*AiSampleCodeResponse, error)
	CreateAiSampleCode(ctx *context.Context, form *api.CreateAiSampleCodesForm) (*discussion_model.AiSampleCode, error)
	DeleteAiSampleCode(ctx *context.Context, id int64) error
}

var DEFAULT_CAPACITY int64 = 10

type SampleCodeServiceImpl struct{}

var _ SampleCodeService = &SampleCodeServiceImpl{}

func (is *SampleCodeServiceImpl) GetAiSampleCodeByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error) {

	response, err := AiSampleCodeDbAdapter.GetAiSampleCodesByCommentID(ctx, commentID)

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


// TODOC 재시도 횟수를 유저 정보로부터 가져와서 제한하기.
// TODOC 잘못된 형식의 json이 돌아올 때 예외 반환하기(json 형식 표시하도록)
func (is *SampleCodeServiceImpl) GenerateAiSampleCodes(ctx *context.Context, form *api.GenerateAiSampleCodesForm) ([]*AiSampleCodeResponse, error) {
	targetCommentId, err := strconv.ParseInt(form.TargetCommentId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid TargetCommentId: %v", err)
	}

	comment, err := issues_model.GetCommentByID(ctx, targetCommentId)
	if err != nil {
		return nil, fmt.Errorf("Comment not found: %v", err)
	}

	issue, err := issues_model.GetIssueByID(ctx, comment.IssueID)
	if err != nil {
		return nil, fmt.Errorf("Issue not found: %v", err)
	}

	repo, err := repo_model.GetRepositoryByID(ctx, issue.RepoID)
	if err != nil {
		return nil, fmt.Errorf("Repository not found: %v", err)
	}

	codeContent, err := GetFileContentFromCommit(ctx, repo.RepoPath(), comment.CommitSHA, comment.TreePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to get file content from commit: %v", err)
	}

	wg := new(sync.WaitGroup)
	wg.Add(AI_SAMPLE_CODE_UNIT)

	resultQueue := make(chan *AiSampleCodeResponse, AI_SAMPLE_CODE_UNIT)

	for i := 0; i < AI_SAMPLE_CODE_UNIT; i++ {
		go func(codeContent, commentContent string) {
			defer wg.Done()

			result, err := AiSampleCodeRequester.RequestReviewToAI(ctx, &AiSampleCodeRequest{
				CodeContent:    codeContent,
				CommentContent: commentContent,
			})

			if err != nil {
				fmt.Errorf("request sample to ai server fail")
				resultQueue <- nil
				return
			}

			highlightedCode, _ := highlight.Code(comment.TreePath, "", result.SampleCode)
			result.SampleCode = string(highlightedCode)

			resultQueue <- result
		}(codeContent, comment.Content)
	}

	wg.Wait()
	close(resultQueue)

	sampleCodes := make([]*AiSampleCodeResponse, 0, AI_SAMPLE_CODE_UNIT)
	for result := range resultQueue {
		sampleCodes = append(sampleCodes, result)
	}

	return sampleCodes, nil
}

// TODOC ai 샘플코드 저장

// TODOC ai 샘플코드 삭제

// TODOC a

func (is *SampleCodeServiceImpl) DeleteAiSampleCode(ctx *context.Context, id int64) error {

	return AiSampleCodeDbAdapter.DeleteAiSampleCodeByID(ctx, id)
}
