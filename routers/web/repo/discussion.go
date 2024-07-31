package repo

import (
	"net/http"

	discussions_model "code.gitea.io/gitea/models/discussions"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	discussions_service "code.gitea.io/gitea/services/discussion"
	"code.gitea.io/gitea/services/forms"
)

const (
	tplDiscussionNew base.TplName = "repo/discussion/new"
)

type RepositoryFile struct {
	Name     string
	NameHash string
	IsBin    bool
}

func NewRepositoryFile(name string) *RepositoryFile {
	repositoryFile := new(RepositoryFile)
	repositoryFile.Name = name
	repositoryFile.NameHash = git.HashFilePathForWebUI(name)
	repositoryFile.IsBin = false
	return repositoryFile
}

// NewDiscussion: render discussion creation page
func NewDiscussion(ctx *context.Context) {

	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true

	// dummy files for repository
	ctx.Data["RepositoryFiles"] = []*RepositoryFile{
		NewRepositoryFile("hello/world.txt"),
		NewRepositoryFile("hello/world.txt"),
		NewRepositoryFile("hello/bye.txt"),
	}
	ctx.HTML(http.StatusOK, tplDiscussionNew)
}

// NewDiscussionPost: serve response for create discussion request
func NewDiscussionPost(ctx *context.Context) {
	repo := ctx.Repo.Repository
	log.Info("NewDiscussionPost has been called")

	log.Info("get form from context")
	form := web.GetForm(ctx).(*forms.CreateDiscussionForm)

	log.Info("form result %v", form)

	ctx.Data["Title"] = ctx.Tr("repo.discussion.new")
	ctx.Data["PageIsDiscussionList"] = true
	log.Info("check ctx written")
	if ctx.Written() {
		return
	}
	log.Info("checked ctx written")

	log.Info("create a new discussion model")
	discussion := &discussions_model.Discussion{
		RepoID:   repo.ID,
		Repo:     repo,
		Title:    form.Title,
		PosterID: ctx.Doer.ID,
		Poster:   ctx.Doer,
		Content:  form.Content,
	}

	log.Info("call new discussion function from discussion service ")
	discussions_service.NewDiscussion(ctx, ctx.Repo.Repository, discussion)
}
