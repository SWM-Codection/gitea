package ai

import (
	"fmt"
	"testing"
	"time"

	"code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAiRequester is a mock implementation of PullCommentRequester
type MockAiRequester struct {
	mock.Mock
}

func (m *MockAiRequester) RequestReviewToAI(ctx *context.Context, request *AiPullCommentRequest) (*AiPullCommentResponse, error) {
	time.Sleep(1000 * time.Millisecond)
	args := m.Called(ctx, request)
	return args.Get(0).(*AiPullCommentResponse), args.Error(1)
}

// MockDbAdapter is a mock implementation of PullCommentDbAdapter
type MockDbAdapter struct {
	mock.Mock
}

func (m *MockDbAdapter) GetIssueByID(ctx *context.Context, id int64) (*issues.Issue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*issues.Issue), args.Error(1)
}

func (m *MockDbAdapter) CreateAiPullComment(ctx *context.Context, opts *issues.CreateAiPullCommentOption) (*issues.AiPullComment, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	commentID := args.Get(0).(*issues.AiPullComment)
	return commentID, args.Error(1)
}

func (m *MockDbAdapter) DeleteAiPullCommentByID(ctx *context.Context, id int64) error {
	args := m.Called(ctx, 1)

	return args.Error(0)
}

func TestCreateAiPullComment(t *testing.T) {
	mockRequester := new(MockAiRequester)
	AiPullCommentRequester = mockRequester
	mockDbAdapter := new(MockDbAdapter)
	AiPullCommentDbAdapter = mockDbAdapter
	aiService := &PullCommentServiceImpl{}
	// Mock context and form
	ctx := &context.Context{}

	fileContent := new([]structs.PathContentMap)
	for i := 0; i < 100; i++ {
		*fileContent = append(*fileContent, structs.PathContentMap{
			TreePath: fmt.Sprintf("file%d.go", i),
			Content:  fmt.Sprintf("code content %d", i),
		})

		mockRequester.On("RequestReviewToAI", ctx, &AiPullCommentRequest{
			Branch:   "main",
			TreePath: fmt.Sprintf("file%d.go", i),
			Content:  fmt.Sprintf("code content %d", i),
		}).Return(&AiPullCommentResponse{
			Branch:   "main",
			TreePath: fmt.Sprintf("file%d.go", i+100),
			Content:  fmt.Sprintf("code content %d", i+100),
		}, nil)

	}

	form := &structs.CreateAiPullCommentForm{
		PullID:       "123",
		Branch:       "main",
		FileContents: fileContent,
	}

	issue := &issues.Issue{}
	mockDbAdapter.On("GetIssueByID", ctx, int64(123)).Return(issue, nil)

	comment := issues.AiPullComment{ID: 10}
	mockDbAdapter.On("CreateAiPullComment", ctx, mock.Anything).Return(&comment, nil)

	err := aiService.CreateAiPullComment(ctx, form)

	assert.NoError(t, err)
	mockRequester.AssertExpectations(t)
	mockDbAdapter.AssertExpectations(t)
}

// TODOC delete 테스트
