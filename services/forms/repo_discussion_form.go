package forms

import (
	"net/http"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/modules/web/middleware"
	"code.gitea.io/gitea/services/context"
	"gitea.com/go-chi/binding"
)

type CreateDiscussionForm struct {
	Name       string
	Content    string
	BranchName string
	Codes      []discussion_client.DiscussionCode
}

func (d *CreateDiscussionForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, d, ctx.Locale)
}

type CreateDiscussionCommentForm struct {
	Content string `form:"content"`

	CodeId    *int64 `form:"codeId"`
	StartLine *int32 `form:"startLine"`
	EndLine   *int32 `form:"endLine"`
}

func (dc *CreateDiscussionCommentForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, dc, ctx.Locale)
}

type ModifyDiscussionCommentForm struct {
	DiscussionCommentId int64  `form:"id"`
	Content             string `form:"content"`

	CodeId    *int64 `form:"codeId"`
	StartLine *int32 `form:"startLine"`
	EndLine   *int32 `form:"endLine"`
}

func (mc *ModifyDiscussionCommentForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, mc, ctx.Locale)
}
