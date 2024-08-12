package ai

import (
	discussion_model "code.gitea.io/gitea/models/discussion"
	"code.gitea.io/gitea/services/context"
)

type DiscussionDbAdapter interface {
	// GetDiscussionCommentByID(ctx *context.Context, id int64) (*issues_model.Issue, error)
	InsertDiscussionAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiSampleCode, error)
	DeleteDiscussionAiSampleCodeByID(ctx *context.Context, id int64) error
}

type DiscussionDbAdapterImpl struct{}

var _ DiscussionDbAdapter = &DiscussionDbAdapterImpl{}

func (is *DiscussionDbAdapterImpl) InsertDiscussionAiSampleCode(ctx *context.Context, opts *discussion_model.CreateDiscussionAiCommentOpt) (*discussion_model.DiscussionAiSampleCode, error) {
	return discussion_model.CreateDiscussionAiSampleCode(ctx, opts)
}

func (is *DiscussionDbAdapterImpl) DeleteDiscussionAiSampleCodeByID(ctx *context.Context, id int64) error {

	return discussion_model.DeleteDiscussionAiSampleCodeByID(ctx, id)

}
