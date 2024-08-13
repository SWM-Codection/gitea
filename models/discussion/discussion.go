package model

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
)

// TODOC 재시도 횟수 저장

type AiSampleCode struct {
	Id              int64 `xorm:"'id' pk autoincr"`
	TargetCommentId int64 `xorm:"'target_comment_id' index notnull"`
	GenearaterId    int64
	Content         string `xorm:"'content'"`
}

type CreateDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64

	Content *string
}

type DeleteDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64
}

var DEFAULT_CAPACITY = 10

func init() {
	db.RegisterModel(new(AiSampleCode))
}

func CreateDiscussionAiSampleCode(ctx context.Context, opts *CreateDiscussionAiCommentOpt) (*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	DiscussionAiComment := &AiSampleCode{
		GenearaterId:    opts.GenearaterId,
		TargetCommentId: opts.TargetCommentId,
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

func DeleteDiscussionAiSampleCodeByID(ctx context.Context, id int64) error {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return err
	}
	defer committer.Close()

	_, err = db.DeleteByID[AiSampleCode](ctx, id)
	if err != nil {
		fmt.Errorf("new Comment insert is invalid")
		return err
	}

	if err = committer.Commit(); err != nil {
		return err
	}

	return nil

}

func GetAiSampleCodeByCommentID(ctx context.Context, id int64) ([]*AiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)

	if err != nil {
		return nil, err
	}

	defer committer.Close()

	e := db.GetEngine(ctx)

	sampleCodes := make([]*AiSampleCode, 0, DEFAULT_CAPACITY)

	err = e.Table("ai_sample_code").Where("comment_id=?", id).Find(&sampleCodes)

	if err != nil {
		return nil, err
	}

	return sampleCodes, nil

}
