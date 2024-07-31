package ai

import (
	"testing"

	"code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAiRequester is a mock implementation of AiRequester
type MockAiRequester struct {
	mock.Mock
}

func (m *MockAiRequester) RequestReviewToAI(ctx *context.Context, request *AiReviewRequest) (*AiReviewResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*AiReviewResponse), args.Error(1)
}

// MockDbAdapter is a mock implementation of DbAdapter
type MockDbAdapter struct {
	mock.Mock
}

func (m *MockDbAdapter) GetIssueByID(ctx *context.Context, id int64) (*issues.Issue, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*issues.Issue), args.Error(1)
}

func (m *MockDbAdapter) CreateAiPullComment(ctx *context.Context, opts *issues.CreateAiPullCommentOption) (*int64, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	commentID := args.Get(0).(*int64)
	return commentID, args.Error(1)
}

func (m *MockDbAdapter) DeleteAiPullCommentByID(ctx *context.Context, id int64) error {
	args := m.Called(ctx, 1)

	return args.Error(0)
}

func TestCreateAiPullComment(t *testing.T) {
	// Set up the mock AiRequester
	mockRequester := new(MockAiRequester)
	aiService := &AiServiceImpl{}

	// Set up the mock DbAdapter
	mockDbAdapter := new(MockDbAdapter)

	// Mock context and form
	ctx := &context.Context{}
	form := &structs.CreateAiPullCommentForm{
		PullID: "123",
		Branch: "main",
		FileContents: &[]structs.PathContentMap{
			{
				TreePath: "file1.go",
				Content:  "code content 1",
			},
			{
				TreePath: "file2.go",
				Content:  "code content 2",
			},
		},
	}

	// Mock response from AI
	mockRequester.On("RequestReviewToAI", ctx, &AiReviewRequest{
		Branch:   "main",
		TreePath: "file1.go",
		Content:  "code content 1",
	}).Return(&AiReviewResponse{
		Branch:   "main",
		TreePath: "file1.go",
		Content:  "reviewed content 1",
	}, nil)

	mockRequester.On("RequestReviewToAI", ctx, &AiReviewRequest{
		Branch:   "main",
		TreePath: "file2.go",
		Content:  "code content 2",
	}).Return(&AiReviewResponse{
		Branch:   "main",
		TreePath: "file2.go",
		Content:  "reviewed content 2",
	}, nil)

	// Mock GetIssueByID
	issue := &issues.Issue{}
	mockDbAdapter.On("GetIssueByID", ctx, int64(123)).Return(issue, nil)

	// Mock CreateAiPullComment
	commentID := int64(1)
	mockDbAdapter.On("CreateAiPullComment", ctx, mock.Anything).Return(&commentID, nil)

	// Call the method under test
	err := aiService.CreateAiPullComment(ctx, form, mockRequester, mockDbAdapter)

	// Assert the expectations
	assert.NoError(t, err)
	mockRequester.AssertExpectations(t)
	mockDbAdapter.AssertExpectations(t)
}

// TODOC delete 테스트
