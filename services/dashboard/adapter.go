package dashboard

import (
	issue_model "code.gitea.io/gitea/models/issues"
	"code.gitea.io/gitea/models/organization"
	repo_model "code.gitea.io/gitea/models/repo"
	"context"
)

type PullDbAdapter interface {
	GetFirstReviewCreatedUnixesByRepoIDs(ctx *context.Context, issueIDs []int64) ([]*issue_model.PullReviewFirstCreatedUnix, error)
	GetOrgRepositoryOrgIDS(ctx *context.Context, orgID int64) ([]int64, error)
}

var _ PullDbAdapter = &PullDbAdapterImpl{}

type PullDbAdapterImpl struct{}

func (is *PullDbAdapterImpl) GetOrgRepositoryOrgIDS(ctx *context.Context, orgID int64) ([]int64, error) {

	return organization.GetOrgRepositoryOrgIDS(*ctx, orgID)
}

func (is *PullDbAdapterImpl) GetOrgRepositories(ctx *context.Context, orgID int64) (repo_model.RepositoryList, error) {
	return organization.GetOrgRepositories(*ctx, orgID)
}

func (is *PullDbAdapterImpl) GetFirstReviewCreatedUnixesByRepoIDs(ctx *context.Context, issueIDs []int64) ([]*issue_model.PullReviewFirstCreatedUnix, error) {
	// 필요한 거 issue의 createdUnix랑 첫 review의 CreatedUnix
	return issue_model.GetFirstReviewCreatedUnixesByRepoIDs(*ctx, issueIDs)
}
