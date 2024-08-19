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

type MockAiSampleCodeRequester struct {
	mock.Mock
}

func (m *MockAiSampleCodeRequester) RequestReviewToAI(ctx *context.Context, request *AiSampleCodeRequest) (*AiSampleCodeResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AiSampleCodeResponse), args.Error(1)
}

// Mock for SampleCodeDbAdapter
type MockDiscussionDbAdapter struct {
	mock.Mock
}

func (m *MockDiscussionDbAdapter) GetAiSampleCodesByCommentID(ctx *context.Context, commentID int64) (*structs.AiSampleCodeResponse, error) {
	args := m.Called(ctx, commentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*structs.AiSampleCodeResponse), args.Error(1)
}

func (m *MockDiscussionDbAdapter) InsertAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.AiSampleCode, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*discussion_model.AiSampleCode), args.Error(1)
}

func (m *MockDiscussionDbAdapter) DeleteAiSampleCodeByID(ctx *context.Context, id int64) error {
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

			mockRequester := new(MockAiSampleCodeRequester)
			AiSampleCodeRequester = mockRequester
			// 모킹을 설정
			service := &SampleCodeServiceImpl{}

			for _, response := range tc.mockResponses {
				mockRequester.On("RequestReviewToAI", mock.Anything, mock.Anything).
					Return(response, tc.mockError).Once()
			}
			ctx := &context.Context{}
			results, err := service.GenerateAiSampleCodes(ctx, tc.form)

			assert.NoError(t, err)
			assert.Len(t, results, tc.expectedLength)

			mockRequester.AssertExpectations(t)
		})
	}
}

func TestCreateAiSampleCode(t *testing.T) {
	tests := []struct {
		name          string
		form          *structs.CreateAiSampleCodesForm
		mockSetup     func(*MockDiscussionDbAdapter)
		expectedCode  *discussion_model.AiSampleCode
		expectedError error
	}{
		{
			name: "Successful creation",
			form: &structs.CreateAiSampleCodesForm{
				TargetCommentId:   "123",
				SampleCodeContent: "Sample code",
			},
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("InsertAiSampleCode", mock.Anything, mock.Anything).Return(&discussion_model.AiSampleCode{Id: 1}, nil)
			},
			expectedCode:  &discussion_model.AiSampleCode{Id: 1},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := new(MockDiscussionDbAdapter)
			AiSampleCodeDbAdapter = mockAdapter
			tt.mockSetup(mockAdapter)

			service := &SampleCodeServiceImpl{}
			ctx := &context.Context{Doer: &user_model.User{ID: 1}}

			code, err := service.CreateAiSampleCode(ctx, tt.form)

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
				m.On("DeleteAiSampleCodeByID", mock.Anything, int64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Not Successful deletion",
			id:   1,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("DeleteAiSampleCodeByID", mock.Anything, int64(1)).Return(errors.New("error"))
			},
			expectedError: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdapter := new(MockDiscussionDbAdapter)
			tt.mockSetup(mockAdapter)

			service := &SampleCodeServiceImpl{}
			ctx := &context.Context{}

			err := service.DeleteAiSampleCode(ctx, tt.id)

			assert.Equal(t, tt.expectedError, err)
			mockAdapter.AssertExpectations(t)
		})
	}
}

func TestGetAiSampleCodeByCommentID(t *testing.T) {
	tests := []struct {
		name          string
		id            int64
		mockSetup     func(*MockDiscussionDbAdapter)
		expectedCode  *structs.AiSampleCodeResponse
		expectedError error
	}{
		{
			name: "success get sample codes with multiple contents",
			id:   1,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				content1 := "example code content 1"
				content2 := "example code content 2"
				m.On("GetAiSampleCodesByCommentID", mock.Anything, int64(1)).Return(&structs.AiSampleCodeResponse{
					CommentID: "1",
					SampleCodeContents: []*structs.AiSampleCodeContent{
						{ID: "3", Content: &content1},
						{ID: "4", Content: &content2},
					},
				}, nil)
			},
			expectedCode: &structs.AiSampleCodeResponse{
				CommentID: "1",
				SampleCodeContents: []*structs.AiSampleCodeContent{
					{ID: "3", Content: func() *string { content := "example code content 1"; return &content }()},
					{ID: "4", Content: func() *string { content := "example code content 2"; return &content }()},
				},
			},
			expectedError: nil,
		},
		{
			name: "error in fetching sample codes",
			id:   2,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("GetAiSampleCodesByCommentID", mock.Anything, int64(2)).Return(nil, errors.New("database error"))
			},
			expectedCode:  nil,
			expectedError: errors.New("database error"),
		},
		{
			name: "no sample codes found",
			id:   3,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				m.On("GetAiSampleCodesByCommentID", mock.Anything, int64(3)).Return(&structs.AiSampleCodeResponse{
					CommentID:          "3",
					SampleCodeContents: []*structs.AiSampleCodeContent{},
				}, nil)
			},
			expectedCode: &structs.AiSampleCodeResponse{
				CommentID:          "3",
				SampleCodeContents: []*structs.AiSampleCodeContent{},
			},
			expectedError: nil,
		},
		{
			name: "success get sample codes with single content",
			id:   4,
			mockSetup: func(m *MockDiscussionDbAdapter) {
				content := "example code content"
				m.On("GetAiSampleCodesByCommentID", mock.Anything, int64(4)).Return(&structs.AiSampleCodeResponse{
					CommentID: "4",
					SampleCodeContents: []*structs.AiSampleCodeContent{
						{ID: "5", Content: &content},
					},
				}, nil)
			},
			expectedCode: &structs.AiSampleCodeResponse{
				CommentID: "4",
				SampleCodeContents: []*structs.AiSampleCodeContent{
					{ID: "5", Content: func() *string { content := "example code content"; return &content }()},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDbAdapter := new(MockDiscussionDbAdapter)
			AiSampleCodeDbAdapter = mockDbAdapter
			tt.mockSetup(mockDbAdapter)

			service := SampleCodeServiceImpl{}
			ctx := &context.Context{}

			result, err := service.GetAiSampleCodeByCommentID(ctx, tt.id)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedCode, result)

			mockDbAdapter.AssertExpectations(t)
		})
	}
}