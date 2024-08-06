package ai

import (
	"errors"
	"testing"

	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	discussion_model "code.gitea.io/gitea/models/discussion"
)

// Mock for AiSampleCodeRequester
type MockAiSampleCodeRequester struct {
	mock.Mock
}

func (m *MockAiSampleCodeRequester) RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*AiSampleCodeResponse), args.Error(1)
}

// Mock for DiscussionDbAdapter
type MockDiscussionDbAdapter struct {
	mock.Mock
}

func (m *MockDiscussionDbAdapter) CreateDiscussionAiComment(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiComment, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*discussion_model.DiscussionAiComment), args.Error(1)
}

func (m *MockDiscussionDbAdapter) DeleteDiscussionAiCommentByID(ctx *context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGenerateAiSampleCodes(t *testing.T) {
	testCases := []struct {
		name           string
		form           *structs.CreateSampleAiCommentForm
		mockResponses  []*AiSampleCodeResponse
		mockError      error
		expectedLength int
		expectedError  bool
	}{
		{
			name: "Successful generation",
			form: &structs.CreateSampleAiCommentForm{
				CodeContent:    "code",
				CommentContent: "comment",
			},
			mockResponses: []*AiSampleCodeResponse{
				{SampleCode: "sample1"},
				{SampleCode: "sample2"},
				{SampleCode: "sample3"},
			},
			mockError:      nil,
			expectedLength: 3,
			expectedError:  false,
		},
		{
			name: "Partial failure",
			form: &structs.CreateSampleAiCommentForm{
				CodeContent:    "code",
				CommentContent: "comment",
			},
			mockResponses: []*AiSampleCodeResponse{
				{SampleCode: "sample1"},
				nil,
				{SampleCode: "sample3"},
			},
			mockError:      errors.New("partial failure"),
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 모킹을 설정
			mockRequester := new(MockAiSampleCodeRequester)
			mockAdapter := new(MockDiscussionDbAdapter)
			service := &DiscussionAiServiceImpl{}

			// 모의 객체의 동작 설정
			for _, response := range tc.mockResponses {
				mockRequester.On("RequestReviewToAI", mock.Anything, mock.Anything).
					Return(response, tc.mockError).Once()
			}
			ctx := &context.Context{}
			// 테스트 대상 함수를 실행
			results, err := service.GenerateAiSampleCodes(ctx, tc.form, mockRequester, mockAdapter)

			assert.NoError(t, err)
			assert.Len(t, results, tc.expectedLength)

			// 모의 객체가 호출 확인
			mockRequester.AssertExpectations(t)
		})
	}
}
