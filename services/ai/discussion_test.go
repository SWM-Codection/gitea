package ai

import (
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

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

func (m *MockDiscussionDbAdapter) InsertDiscussionAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiSampleCode, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*discussion_model.DiscussionAiSampleCode), args.Error(1)
}

func (m *MockDiscussionDbAdapter) DeleteDiscussionAiSampleCodeByID(ctx *context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGenerateAiSampleCodes(t *testing.T) {
	testCases := []struct {
		name           string
		form           *structs.GenerateAiSampleCodesForm
		mockResponses  []*AiSampleCodeResponse
		mockError      error
		expectedLength int
		expectedError  bool
	}{
		{
			name: "Successful generation",
			form: &structs.GenerateAiSampleCodesForm{
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
			form: &structs.GenerateAiSampleCodesForm{
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

func TestCreateAiSampleCode(t *testing.T) {
	tests := []struct {
		name          string
		form          *structs.CreateAiSampleCodesForm
		mockSetup     func(*MockDiscussionDbAdapter)
		expectedCode  *discussion_model.DiscussionAiSampleCode
		expectedError error
	}{
		{
			name: "Successful creation",
			form: &structs.CreateAiSampleCodesForm{
				TargetCommentId:   "123",
				SampleCodeContent: "Sample code",
			},
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("InsertDiscussionAiSampleCode", mock.Anything, mock.Anything).Return(&discussion_model.DiscussionAiSampleCode{Id: 1}, nil)
			},
			expectedCode:  &discussion_model.DiscussionAiSampleCode{Id: 1},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := new(MockDiscussionDbAdapter)
			tt.mockSetup(mockAdapter)

			service := &DiscussionAiServiceImpl{}
			ctx := &context.Context{Doer: &user_model.User{ID: 1}}

			code, err := service.CreateAiSampleCode(ctx, tt.form, mockAdapter)

			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedError, err)
			mockAdapter.AssertExpectations(t)
		})
	}
}

func TestDeleteAiSampleCode(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*MockDiscussionDbAdapter)
		expectedError error
	}{
		{
			name: "Successful deletion",
			id:   1,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("DeleteDiscussionAiSampleCodeByID", mock.Anything, int64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Not Successful deletion",
			id:   1,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("DeleteDiscussionAiSampleCodeByID", mock.Anything, int64(1)).Return(errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := new(MockDiscussionDbAdapter)
			tt.mockSetup(mockAdapter)

			service := &DiscussionAiServiceImpl{}
			ctx := &context.Context{}

			err := service.DeleteAiSampleCode(ctx, tt.id, mockAdapter)

			assert.Equal(t, tt.expectedError, err)
			mockAdapter.AssertExpectations(t)
		})
	}
}
