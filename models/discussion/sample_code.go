package discussion

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/modules/timeutil"

	"code.gitea.io/gitea/models/db"
)

func init() {
	db.RegisterModel(new(AiSampleCode))
}

// TODOC 재시도 횟수 저장

type AiSampleCode struct {
	TargetCommentId int64              `xorm:"'target_comment_id'"`
	Id              int64              `xorm:"'id' pk autoincr"`
	CodeId          int64              `xorm:"'code_id'"`
	StartLine       int64              `xorm:"'start_line'"`
	EndLine         int64              `xorm:"'end_line'"`
	DiscussionId    int64              `xorm:"'discussion_id'"`
	GenearaterId    int64              `xorm:"'genearater_id'"`
	CommentType     string             `xorm:"'comment_type'"`
	Content         string             `xorm:"'content' text"`
	CreatedUnix     timeutil.TimeStamp `xorm:"INDEX created"`
	UpdatedUnix     timeutil.TimeStamp `xorm:"INDEX updated"`
	Status          string             `xorm:"status"` //
	DeletedUnix     timeutil.TimeStamp `xorm:"deleted"`
}

type CreateDiscussionAiCommentOpt struct {
	CodeId       int64
	StartLine    int64
	EndLine      int64
	GenearaterId int64
	DiscussionId int64
	Type         string
	Content      *string
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
		GenearaterId: opts.GenearaterId,
		CodeId:       opts.CodeId,
		DiscussionId: opts.DiscussionId,
		StartLine:    opts.StartLine,
		EndLine:      opts.EndLine,
		CommentType:  opts.Type,
		Content:      *opts.Content,
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

func GetAiSampleCodesByCodeId(ctx context.Context, id int64, sampleType string) ([]*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)

	if err != nil {
		return nil, err
	}

	defer committer.Close()

	e := db.GetEngine(ctx)

	var sampleCode []*AiSampleCode

	err = e.Table("ai_sample_code").Where("code_id=?", id).Find(&sampleCode)

	if err != nil {
		return nil, fmt.Errorf("failed to get AI sample code: %v", err)
	}

	return sampleCode, nil
}

func GetAiSampleCodeById(ctx context.Context, id int64) (*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	var sampleCode *AiSampleCode

	e := db.GetEngine(ctx)
	has, err := e.Table("ai_sample_code").Where("codeId=?", id).Get(sampleCode)

	if !has {
		return nil, nil
	}

	return sampleCode, nil
}
