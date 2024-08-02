package forms

import (
	"net/http"

	discussion_client "code.gitea.io/gitea/client/discussion"
	"code.gitea.io/gitea/modules/web/middleware"
	"code.gitea.io/gitea/services/context"
	"gitea.com/go-chi/binding"
)

type CreateDiscussionForm struct {
	Name    string
	Content string
	Codes   []discussion_client.DiscussionCode
}

func (d *CreateDiscussionForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	ctx := context.GetValidateContext(req)
	return middleware.Validate(errs, ctx.Data, d, ctx.Locale)
}
