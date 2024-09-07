package discussion

import (
	"code.gitea.io/gitea/modules/timeutil"
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
)

func init() {
	db.RegisterModel(new(AiSampleCode))
}

// TODOC 재시도 횟수 저장

type AiSampleCode struct {
	Id              int64 `xorm:"'id' pk autoincr"`
	TargetCommentId int64 `xorm:"'target_comment_id' INDEX NOT NULL"`
	GenearaterId    int64
	CommentType			string			   `xorm:"'comment_type'"`
	Content         string             `xorm:"'content'"`
	CreatedUnix     timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix     timeutil.TimeStamp `xorm:"INDEX updated"`
	Status          string             `xorm:"status"` //
	DeletedUnix     timeutil.TimeStamp `xorm:"deleted"`
}

type CreateDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64
	Type			string

	Content *string
}

type DeleteDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64
}

var DEFAULT_CAPACITY = 10

func CreateAiSampleCode(ctx context.Context, opts *CreateDiscussionAiCommentOpt) (*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	DiscussionAiComment := &AiSampleCode{
		GenearaterId:    opts.GenearaterId,
		TargetCommentId: opts.TargetCommentId,
		CommentType: opts.Type,
		Content:         *opts.Content,
	}

	e := db.GetEngine(ctx)
	_, err = e.Insert(DiscussionAiComment)
	if err != nil {
		fmt.Errorf("new Comment insert is invalid")
		return nil, err
	}

	if err = committer.Commit(); err != nil {
		return nil, err
	}

	return DiscussionAiComment, nil

}

func DeleteAiSampleCodeByID(ctx context.Context, id int64) error {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return err
	}
	defer committer.Close()

	_, err = db.DeleteByID[AiSampleCode](ctx, id)
	if err != nil {
		fmt.Errorf("sample code not exists")
		return err
	}

	if err = committer.Commit(); err != nil {
		return err
	}

	return nil

}

func GetAiSampleCodeByCommentID(ctx context.Context, id int64, sampleType string) ([]*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)

	if err != nil {
		return nil, err
	}

	defer committer.Close()

	e := db.GetEngine(ctx)

	sampleCodes := make([]*AiSampleCode, 0, DEFAULT_CAPACITY)

	err = e.Table("ai_sample_code").Where("target_comment_id=? AND comment_type = ?", id, sampleType).Find(&sampleCodes)

	if err != nil {
		fmt.Errorf(err.Error())
		return nil, err
	}

	return sampleCodes, nil
}
