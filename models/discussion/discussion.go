package model

import (
	"context"
	"fmt"

	"code.gitea.io/gitea/models/db"
)

// TODOC 재시도 횟수 저장

type AiSampleCodeRetryCount struct {
	TargetCommentId int64 `xorm: "target_comment_id"`
	count           int64 `xorm: "count"`
}

type DiscussionAiComment struct {
	Id              int64 `xorm:"'id' pk autoincr"`
	TargetCommentId int64 `xorm:"'discussion_id' notnull"`
	GenearaterId    int64
	Content         string `xorm:"'content'"`
}

type CreateDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64
	// Scope        string `xorm:"'scope'"`
	Content string
}

type DeleteDiscussionAiCommentOpt struct {
	TargetCommentId int64
	GenearaterId    int64
	// Scope        string `xorm:"'scope'"`
}

type CommentScopeEnum string

const (
	GLOBAL CommentScopeEnum = "GLOBAL"
	LOCAL  CommentScopeEnum = "LOCAL"
)

func init() {
	db.RegisterModel(new(DiscussionAiComment))
}

func CreateDiscussionAiComment(ctx context.Context, opts *CreateDiscussionAiCommentOpt) (*DiscussionAiComment, error) {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return nil, err
	}
	defer committer.Close()

	DiscussionAiComment := &DiscussionAiComment{
		GenearaterId:    opts.GenearaterId,
		TargetCommentId: opts.TargetCommentId,
		Content:         opts.Content,
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

func DeleteDiscussionAiCommentByID(ctx context.Context, id int64) error {

	ctx, committer, err := db.TxContext(ctx)
	if err != nil {
		return err
	}
	defer committer.Close()

	_, err = db.DeleteByID[DiscussionAiComment](ctx, id)
	if err != nil {
		fmt.Errorf("new Comment insert is invalid")
		return err
	}

	if err = committer.Commit(); err != nil {
		return err
	}

	return nil

}
