package ai

import (
	discussion_model "code.gitea.io/gitea/models/discussion"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/context"
	"github.com/spf13/cast"
)

type SampleCodeDbAdapter interface {
	// GetDiscussionCommentByID(ctx *context.Context, id int64) (*issues_model.Issue, error)
	InsertAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.AiSampleCode, error)
	DeleteAiSampleCodeByID(ctx *context.Context, id int64) error
	GetAiSampleCodesByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error)
}

type SampleCodeDbAdapterImpl struct{}

var _ SampleCodeDbAdapter = &SampleCodeDbAdapterImpl{}

func (is *SampleCodeDbAdapterImpl) GetAiSampleCodesByCommentID(ctx *context.Context, commentID int64) (*api.AiSampleCodeResponse, error) {

	sampleCodes, err := discussion_model.GetAiSampleCodeByCommentID(ctx, commentID)

	if err != nil {
		return nil, err
	}

	response := api.AiSampleCodeResponse{
		CommentID:          cast.ToString(commentID),
		SampleCodeContents: make([]*api.AiSampleCodeContent, 0, DEFAULT_CAPACITY),
	}

	for _, sampleCode := range sampleCodes {
		response.SampleCodeContents = append(response.SampleCodeContents, &api.AiSampleCodeContent{
			ID:      cast.ToString(sampleCode.Id),
			Content: &sampleCode.Content,
		})
	}

	return &response, nil
}

func (is *SampleCodeDbAdapterImpl) InsertAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.AiSampleCode, error) {
	return discussion_model.CreateAiSampleCode(ctx, opts)
}

func (is *SampleCodeDbAdapterImpl) DeleteAiSampleCodeByID(ctx *context.Context, id int64) error {

	return discussion_model.DeleteAiSampleCodeByID(ctx, id)

}
