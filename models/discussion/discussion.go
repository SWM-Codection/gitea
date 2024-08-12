package model

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
)

// TODOC 재시도 횟수 저장

type DiscussionAiSampleCode struct {
	Id              int64 `xorm:"'id' pk autoincr"`
	TargetCommentId int64 `xorm:"'discussion_id' notnull"`
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

func init() {
	db.RegisterModel(new(DiscussionAiSampleCode))
}

func CreateDiscussionAiSampleCode(ctx context.Context, opts *CreateDiscussionAiCommentOpt) (*DiscussionAiSampleCode, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	DiscussionAiComment := &DiscussionAiSampleCode{
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

	_, err = db.DeleteByID[DiscussionAiSampleCode](ctx, id)
	if err != nil {
		fmt.Errorf("new Comment insert is invalid")
		return err
	}

	if err = committer.Commit(); err != nil {
		return err
	}

	return nil

}
