// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues

import (
	"context"

	"code.gitea.io/gitea/models/db"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/markup"
	"code.gitea.io/gitea/modules/markup/markdown"
	discussion_model "code.gitea.io/gitea/models/discussion"

	"xorm.io/builder"
)

// CodeComments represents comments on code by using this structure: FILENAME -> LINE (+ == proposed; - == previous) -> COMMENTS
type CodeComments map[string]map[int64][]*Comment

// FetchCodeComments will return a 2d-map: ["Path"]["Line"] = Comments at line
func FetchCodeComments(ctx context.Context, issue *Issue, currentUser *user_model.User, showOutdatedComments bool) (CodeComments, error) {
	return fetchCodeCommentsByReview(ctx, issue, currentUser, nil, showOutdatedComments)
}

func fetchCodeCommentsByReview(ctx context.Context, issue *Issue, currentUser *user_model.User, review *Review, showOutdatedComments bool) (CodeComments, error) {
	pathToLineToComment := make(CodeComments)
	if review == nil {
		review = &Review{ID: 0}
	}
	opts := FindCommentsOptions{
		Type:     CommentTypeCode,
		IssueID:  issue.ID,
		ReviewID: review.ID,
	}

	comments, err := findCodeComments(ctx, opts, issue, currentUser, review, showOutdatedComments)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		if pathToLineToComment[comment.TreePath] == nil {
			pathToLineToComment[comment.TreePath] = make(map[int64][]*Comment)
		}
		pathToLineToComment[comment.TreePath][comment.Line] = append(pathToLineToComment[comment.TreePath][comment.Line], comment)

		aiSampleCode, err := discussion_model.GetAiSampleCodeByCommentID(ctx, comment.ID, "pull")
		if err != nil {
			return nil, err
		}

		if aiSampleCode != nil {
			aiComment, err := convertAiSampleCodeToComment(
				ctx,
				aiSampleCode,
				issue,
				comment,
			)
			if err != nil {
				return nil, err
			}

			pathToLineToComment[comment.TreePath][comment.Line] = append(pathToLineToComment[comment.TreePath][comment.Line], aiComment)
		}
	}
	
	return pathToLineToComment, nil
}

func findCodeComments(ctx context.Context, opts FindCommentsOptions, issue *Issue, currentUser *user_model.User, review *Review, showOutdatedComments bool) ([]*Comment, error) {
	var comments CommentList
	if review == nil {
		review = &Review{ID: 0}
	}
	conds := opts.ToConds()

	if !showOutdatedComments && review.ID == 0 {
		conds = conds.And(builder.Eq{"invalidated": false})
	}

	e := db.GetEngine(ctx)
	if err := e.Where(conds).
		Asc("comment.created_unix").
		Asc("comment.id").
		Find(&comments); err != nil {
		return nil, err
	}

	if err := issue.LoadRepo(ctx); err != nil {
		return nil, err
	}

	if err := comments.LoadPosters(ctx); err != nil {
		return nil, err
	}

	if err := comments.LoadAttachments(ctx); err != nil {
		return nil, err
	}

	// Find all reviews by ReviewID
	reviews := make(map[int64]*Review)
	ids := make([]int64, 0, len(comments))
	for _, comment := range comments {
		if comment.ReviewID != 0 {
			ids = append(ids, comment.ReviewID)
		}
	}
	if err := e.In("id", ids).Find(&reviews); err != nil {
		return nil, err
	}

	n := 0
	for _, comment := range comments {
		if re, ok := reviews[comment.ReviewID]; ok && re != nil {
			// If the review is pending only the author can see the comments (except if the review is set)
			if review.ID == 0 && re.Type == ReviewTypePending &&
				(currentUser == nil || currentUser.ID != re.ReviewerID) {
				continue
			}
			comment.Review = re
		}
		comments[n] = comment
		n++

		if err := comment.LoadResolveDoer(ctx); err != nil {
			return nil, err
		}

		if err := comment.LoadReactions(ctx, issue.Repo); err != nil {
			return nil, err
		}

		var err error
		if comment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
			Ctx: ctx,
			Links: markup.Links{
				Base: issue.Repo.Link(),
			},
			Metas: issue.Repo.ComposeMetas(ctx),
		}, comment.Content); err != nil {
			return nil, err
		}
	}
	return comments[:n], nil
}

// FetchCodeCommentsByLine fetches the code comments for a given treePath and line number
func FetchCodeCommentsByLine(ctx context.Context, issue *Issue, currentUser *user_model.User, treePath string, line int64, showOutdatedComments bool) (CommentList, error) {
	opts := FindCommentsOptions{
		Type:     CommentTypeCode,
		IssueID:  issue.ID,
		TreePath: treePath,
		Line:     line,
	}
	return findCodeComments(ctx, opts, issue, currentUser, nil, showOutdatedComments)
}

func FetchCodeAiComments(ctx context.Context, issue *Issue, fileLines map[string]int64) (CodeComments, error) {
	return fetchCodeAiCommentsByReview(ctx, issue, fileLines)
}

func fetchCodeAiCommentsByReview(ctx context.Context, issue *Issue, fileLines map[string]int64) (CodeComments, error) {
	pathToLineToComment := make(CodeComments)
	var comments CommentList

	// AiPullComment 리스트를 가져옴
	aiPullComments, err := fetchAiPullComments(ctx, issue)
	if err != nil {
		return nil, err
	}

	// AiPullComment를 Comment로 변환
	for _, aiPullComment := range aiPullComments {
		line := fileLines[aiPullComment.TreePath]
		updateOpts := UpdateAiPullCommentOption{
			CommentID:	aiPullComment.ID,
			Line:		&line,
		}
		updatedAiPullComment, err := UpdateAiPullComment(ctx, &updateOpts)
		if err != nil {
			return nil, err
		}
		comment, err := convertAiPullCommentToComment(ctx, updatedAiPullComment, issue)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	for _, comment := range comments {
		if pathToLineToComment[comment.TreePath] == nil {
			pathToLineToComment[comment.TreePath] = make(map[int64][]*Comment)
		}
		pathToLineToComment[comment.TreePath][comment.Line] = append(pathToLineToComment[comment.TreePath][comment.Line], comment)
	}

	return pathToLineToComment, nil
}

func FetchAiPullCommentByLine(ctx context.Context, issue *Issue, treePath string, line int64) (*Comment, error) {
	aiPullComment, err := fetchAiPullCommentByLine(ctx, issue, treePath, line)
	if err != nil || aiPullComment == nil {
		return nil, err
	}
	return convertAiPullCommentToComment(ctx, aiPullComment, issue)
}

// convertAiPullCommentToComment converts an AiPullComment into a Comment
func convertAiPullCommentToComment(ctx context.Context, aiPullComment *AiPullComment, issue *Issue) (*Comment, error) {
	// AiPullComment를 Comment로 변환
	comment := &Comment{
		ID:          0,
		PosterID:    -3,
		IssueID:     aiPullComment.PullID,
		Content:     aiPullComment.Content,
		TreePath:    aiPullComment.TreePath,
		CreatedUnix: aiPullComment.CreatedUnix,
		UpdatedUnix: aiPullComment.UpdatedUnix,
		CommitSHA:   aiPullComment.CommitSHA,
		Line:        aiPullComment.Line,
		//Poster:		 user_model.NewCodectionUser(),
	}

	if err := comment.LoadPoster(ctx); err != nil {
		return nil, err
	}

	if err := comment.LoadAttachments(ctx); err != nil {
		return nil, err
	}

	if err := comment.LoadReactions(ctx, issue.Repo); err != nil {
		return nil, err
	}

	// 마크다운 렌더링 (AI Pull Comment의 내용을 HTML로 변환)
	var err error
	if comment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
		Ctx: ctx,
		Links: markup.Links{
			Base: issue.Repo.Link(),
		},
		Metas: issue.Repo.ComposeMetas(ctx),
	}, aiPullComment.Content); err != nil {
		return nil, err
	}

	return comment, nil
}

func MergeAIComments(allComments, aiComments CodeComments) {
	for fileName, aiLineCommits := range aiComments {
		if existingLineCommits, ok := allComments[fileName]; ok {
			for lineNumber, aiCommentsForLine := range aiLineCommits {
				if existingCommentsForLine, exists := existingLineCommits[lineNumber]; exists {
					allComments[fileName][lineNumber] = append(existingCommentsForLine, aiCommentsForLine...)
				} else {
					allComments[fileName][lineNumber] = aiCommentsForLine
				}
			}
		} else {
			allComments[fileName] = aiLineCommits
		}
	}
}

func convertAiSampleCodeToComment(ctx context.Context, aiSampleCode *discussion_model.AiSampleCode, issue *Issue, target_comment *Comment) (*Comment, error) {
	comment := &Comment{
		ID:          -target_comment.ID,
		PosterID:    -3,
		IssueID:     aiSampleCode.TargetCommentId,
		Content:     aiSampleCode.Content,
		CreatedUnix: aiSampleCode.CreatedUnix,
		UpdatedUnix: aiSampleCode.UpdatedUnix,
		CommitSHA:   target_comment.CommitSHA,
		Line:        target_comment.Line,
	}

	if err := comment.LoadPoster(ctx); err != nil {
		return nil, err
	}

	if err := comment.LoadAttachments(ctx); err != nil {
		return nil, err
	}

	if err := comment.LoadReactions(ctx, issue.Repo); err != nil {
		return nil, err
	}

	var err error
	if comment.RenderedContent, err = markdown.RenderString(&markup.RenderContext{
		Ctx: ctx,
		Links: markup.Links{
			Base: issue.Repo.Link(),
		},
		Metas: issue.Repo.ComposeMetas(ctx),
	}, aiSampleCode.Content); err != nil {
		return nil, err
	}

	return comment, nil
}