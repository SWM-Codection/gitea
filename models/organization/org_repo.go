// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package organization

import (
	"context"

	"code.gitea.io/gitea/models/db"
	repo_model "code.gitea.io/gitea/models/repo"
)

// GetOrgRepositories get repos belonging to the given organization
func GetOrgRepositories(ctx context.Context, orgID int64) (repo_model.RepositoryList, error) {
	var orgRepos []*repo_model.Repository
	return orgRepos, db.GetEngine(ctx).Where("owner_id = ?", orgID).Find(&orgRepos)
}

func GetOrgRepositoryOrgIDS(ctx context.Context, orgID int64) ([]int64, error) {
	var orgRepos []int64
	err := db.GetEngine(ctx).Table("repository").
		Where("owner_id = ?", orgID).
		Cols("id").
		Find(&orgRepos)
	return orgRepos, err
}
