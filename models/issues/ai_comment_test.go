package issues_test

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	ai_comment "code.gitea.io/gitea/models/issues"
)

func TestCreateAiPullComment(t *testing.T) {

	assert := assert.New(t)
	assert.NoError(unittest.PrepareTestDatabase())

	pull := unittest.AssertExistsAndLoadBean(t, &ai_comment.Issue{ID: 2})
	doer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: pull.RepoID})
	content := "this code is sh**"
	treePath := "/src/ddd"
	newAiCommentID, err := ai_comment.CreateAiPullComment(db.DefaultContext, &ai_comment.CreateAiPullCommentOption{
		Doer:     doer,
		Pull:     pull,
		Repo:     repo,
		Content:  content,
		TreePath: treePath,
	})
	assert.NoError(err)
	newAiComment, err := ai_comment.GetAIPullCommentByID(db.DefaultContext, *newAiCommentID)
	if err != nil {
		return
	}

	assert.EqualValues(newAiComment.Content, content)
	assert.EqualValues(newAiComment.TreePath, treePath)
	assert.EqualValues(newAiComment.PosterID, doer.ID)
	assert.EqualValues(newAiComment.PullID, pull.ID)

}

func TestDeleteAiPullRequest(t *testing.T) {


	assert.NoError(t, unittest.PrepareTestDatabase())

	comment := unittest.AssertExistsAndLoadBean(t, &ai_comment.AiPullComment{ID: 2})
	assert.NoError(t, ai_comment.DeleteAiPullCommentByID(db.DefaultContext, comment.ID))
	unittest.AssertNotExistsBean(t, &ai_comment.AiPullComment{ID: comment.ID})

	assert.NoError(t, ai_comment.DeleteAiPullCommentByID(db.DefaultContext, unittest.NonexistentID))
	unittest.CheckConsistencyFor(t, &ai_comment.AiPullComment{})
}
