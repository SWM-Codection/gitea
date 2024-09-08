package issues

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"

	"code.gitea.io/gitea/modules/timeutil"
)

// init 메소드가 있으면 자동적으로 xorm에서 이 메소드를 실행하는듯 하다.
func init() {
	db.RegisterModel(new(AiPullComment))
}

// TODOC AI 코멘트 테이블 만들기
// TODOC outdated가 어떤 식으로 나타나는 것인지 알아보기
// TODOC 먼저 영속성 계층부터-도메인 계층 순서로 만들어가기
type AiPullComment struct {
	ID          int64            `xorm:"pk autoincr"`
	PosterID    int64            `xorm:"INDEX"`
	Poster      *user_model.User `xorm:"-"`
	PullID      int64
	Pull        *Issue `xorm:"-"`
	TreePath    string
	Content     string             `xorm:"LONGTEXT"`
	CreatedUnix timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix timeutil.TimeStamp `xorm:"INDEX updated"`
	Status      string             `xorm:"status"` //
	DeletedUnix timeutil.TimeStamp `xorm:"deleted"`
	CommitSHA   string             `xorm:"VARCHAR(64)"`
	Line        int64
	// CommitID        int64
}

type FindAiPullCommentsOptions struct {
	RepoID    int64
	PullID    int64
}

type CreateAiPullCommentOption struct {
	Doer      *user_model.User
	Repo      *repo_model.Repository
	Pull      *Issue
	TreePath  string
	Content   string
	CommitSHA string
	// CommitID string

}

type UpdateAiPullCommentOption struct {
	CommentID  int64
	Line	   *int64
	Content	   *string
}

type ErrAiPullCommentNotExist struct {
	ID     int64
	RepoID int64
}

func IsErrAiPullCommentNotExist(err error) bool {
	_, ok := err.(ErrIssueWasClosed)
	return ok
}

func (err ErrAiPullCommentNotExist) Error() string {
	return fmt.Sprintf("AiPullComment does not exist [id: %d, repo_id: %d]", err.ID, err.RepoID)
}

func CreateAiPullComment(ctx context.Context, opts *CreateAiPullCommentOption) (*AiPullComment, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	aiPullComment := &AiPullComment{
		PosterID:  opts.Doer.ID,
		PullID:    opts.Pull.ID,
		TreePath:  opts.TreePath,
		Content:   opts.Content,
		CommitSHA: opts.CommitSHA,
	}

	e := db.GetEngine(ctx)
	_, err = e.Insert(aiPullComment)
	if err != nil {
		fmt.Errorf("new Comment insert is invalid")
		return nil, err
	}

	if err = committer.Commit(); err != nil {
		return nil, err
	}

	return aiPullComment, nil

}

func GetAIPullCommentByID(ctx context.Context, id int64) (*AiPullComment, error) {
	comment := new(AiPullComment)
	has, err := db.GetEngine(ctx).ID(id).Get(comment)

	if err != nil {
		return nil, ErrAiPullCommentNotExist{id, 0}
	} else if !has {
		return nil, err
	}

	return comment, nil
}

func DeleteAiPullCommentByID(ctx context.Context, id int64) error {
	_, err := GetAIPullCommentByID(ctx, id)

	if err != nil {

		if IsErrAiPullCommentNotExist(err) {
			return nil
		}
		return err

	}

	dbCtx, commiter, err := db.TxContext(ctx)

	defer commiter.Close()

	if err != nil {
		return err
	}

	_, err = db.DeleteByID[AiPullComment](dbCtx, id)

	if err != nil {
		return err
	}

	return commiter.Commit()

}

func UpdateAiPullComment(ctx context.Context, opts *UpdateAiPullCommentOption) (*AiPullComment, error) {
	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	aiPullComment := &AiPullComment{}
	e := db.GetEngine(ctx)
	has, err := e.ID(opts.CommentID).Get(aiPullComment)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("ai pull comment not found")
	}

	// Update the fields only if the new values are provided
	if opts.Line != nil {
		aiPullComment.Line = *opts.Line
	}
	if opts.Content != nil {
		aiPullComment.Content = *opts.Content
	}

	aiPullComment.UpdatedUnix = timeutil.TimeStampNow()

	if _, err := e.ID(aiPullComment.ID).AllCols().Update(aiPullComment); err != nil {
		return nil, err
	}

	if err := committer.Commit(); err != nil {
		return nil, err
	}

	return aiPullComment, nil
}

type AiPullCommentList []*AiPullComment

func fetchAiPullComments(ctx context.Context, issue *Issue) ([]*AiPullComment, error) {
	var aiPullComments []*AiPullComment

	e := db.GetEngine(ctx)
	if err := e.Where("pull_id = ?", issue.ID).Find(&aiPullComments); err != nil {
		return nil, err
	}

	return aiPullComments, nil
}

func fetchAiPullCommentByLine(ctx context.Context, issue *Issue, treePath string, line int64) (*AiPullComment, error) {
	var aiPullComment AiPullComment

	e := db.GetEngine(ctx)
	if has, err := e.Where("pull_id = ? AND tree_path = ? AND line = ?", issue.ID, treePath, line).Get(&aiPullComment); err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}

	return &aiPullComment, nil
}


// TODOC repo가 삭제되면 Ai Comment도 삭제하는 로직
