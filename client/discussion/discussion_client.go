package discussion

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"code.gitea.io/gitea/services/context"

	"code.gitea.io/gitea/client"
	repo_model "code.gitea.io/gitea/models/repo"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/timeutil"
	"github.com/go-resty/resty/v2"
)

type DiscussionCode struct {
	Id        int64  `json:"id"`
	FilePath  string `json:"filePath"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}

type DeleteDiscussionCommentRequest struct {
	PosterId            int64 `json:"posterId"`
	DiscussionCommentId int64 `json:"discussionCommentId"`
}

type PostDiscussionRequest struct {
	RepoId     int64            `json:"repoId"`
	Poster     *user_model.User `json:"-"`
	PosterId   int64            `json:"posterId"`
	Name       string           `json:"name"`
	Content    string           `json:"content"`
	BranchName string           `json:"branchName"`
	Codes      []DiscussionCode `json:"codes"`
}

type DiscussionAvailableRequest struct {
	RepoId    int64 `json:"repoId"`
	Available bool  `json:"available"`
}

type CommentScopeEnum int

const (
	CommentScopeGlobal CommentScopeEnum = iota
	CommentScopeLocal
)

type PostCommentRequest struct {
	DiscussionId int64            `json:"discussionId"`
	CodeId       *int64           `json:"codeId"`
	PosterId     int64            `json:"posterId"`
	Scope        CommentScopeEnum `json:"scope"`
	StartLine    *int32           `json:"startLine"`
	EndLine      *int32           `json:"endLine"`
	Content      string           `json:"content"`
}

type ModifyDiscussionRequest struct {
	RepoId       int64            `json:"repoId"`
	DiscussionId int64            `json:"discussionId"`
	PosterId     int64            `json:"posterId"`
	Name         string           `json:"name"`
	Content      string           `json:"content"`
	Codes        []DiscussionCode `json:"codes"`
}

type DiscussionResponse struct {
	Id          int64              `json:"id"`
	Name        string             `json:"name"`
	Content     string             `json:"content"`
	RepoId      int64              `json:"repoId"`
	PosterId    int64              `json:"posterId"`
	CommitHash  string             `json:"commitHash"`
	IsClosed    bool               `json:"isClosed"`
	Deadline    timeutil.TimeStamp `json:"deadline"`
	Assignees   []int64            `json:"assignees"`
	CreatedUnix timeutil.TimeStamp `json:"createdUnix"`
	UpdatedUnix timeutil.TimeStamp `json:"updatedUnix"`
	Index       int64              `json:"index"`
	Poster      *user_model.User   `json:"-"`
}

type Discussion struct {
	Id           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Content      string                 `json:"content"`
	RepoId       int64                  `json:"repoId"`
	PosterId     int64                  `json:"posterId"`
	CommitHash   string                 `json:"commitHash"`
	Index        int64                  `json:"index"`
	IsClosed     bool                   `json:"isClosed"`
	CreatedUnix  timeutil.TimeStamp     `json:"createdUnix"` // required, but didn't exist before
	ClosedUnix   timeutil.TimeStamp     `json:"closedUnix"`  // required, but didn't exist before
	DeadlineUnix timeutil.TimeStamp     `json:"deadlineUnix"`
	NumComments  int                    `json:"-"` // it can be computed
	Repo         *repo_model.Repository `json:"-"` // it can be computed via d.LoadRepo()
	Poster       *user_model.User       `json:"-"` // it can be computed via d.LoadPoster()
}

type DiscussionListResponse struct {
	TotalCount  int64         `json:"totalCount"`
	Discussions []*Discussion `json:"discussions"`
}

type ModifyDiscussionCommentRequest struct {
	DiscussionId        int64            `json:"discussionId"`
	DiscussionCommentId int64            `json:"discussionCommentId"`
	CodeId              *int64           `json:"codeId"`
	PosterId            int64            `json:"posterId"`
	Scope               CommentScopeEnum `json:"scope"`
	StartLine           *int32           `json:"startLine"`
	EndLine             *int32           `json:"endLine"`
	Content             string           `json:"content"`
}

type DiscussionCountResponse struct {
	OpenCount  int `json:"openCount"`
	CloseCount int `json:"closeCount"`
}

type ExtractedLine struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
}
type CodeBlock struct {
	CodeId   int64                       `json:"codeId"`
	Lines    []ExtractedLine             `json:"lines"`
	Comments []DiscussionCommentResponse `json:"comments"`
}
type FileContent struct {
	FilePath   string      `json:"filePath"`
	CodeBlocks []CodeBlock `json:"codeBlocks"`
}

type DiscussionCommentResponse struct {
	Id           int64                 `json:"id"`
	DiscussionId int64                 `json:"discussionId"`
	PosterId     int64                 `json:"posterId"`
	CodeId       int64                 `json:"codeId"`
	Scope        string                `json:"scope"`
	StartLine    int64                 `json:"startLine"`
	EndLine      int64                 `json:"endLine"`
	Content      string                `json:"content"`
	CreatedUnix  timeutil.TimeStamp    `json:"createdUnix"`
	Reactions    []*DiscussionReaction `json:"reactions"`
}

type ReactionTypeEnum = string

const (
	PLUS_ONE  ReactionTypeEnum = "+1"
	MINUS_ONE ReactionTypeEnum = "-1"
	LAUGH     ReactionTypeEnum = "laugh"
	HOORAY    ReactionTypeEnum = "hooray"
	CONFUSED  ReactionTypeEnum = "confused"
	HEART     ReactionTypeEnum = "heart"
	ROCKET    ReactionTypeEnum = "rocket"
	EYES      ReactionTypeEnum = "eyes"
)

type DiscussionReaction struct {
	Id           int64            `json:"id"`
	Type         ReactionTypeEnum `json:"type"`
	DiscussionId int64            `json:"discussionId"`
	CommentId    int64            `json:"commentId"`
	UserId       int64            `json:"userId"`
}
type DiscussionContentResponse struct {
	DiscussionId    int64                       `json:"discussionId"`
	Contents        []FileContent               `json:"contents"`
	GlobalComments  []DiscussionCommentResponse `json:"globalComments"`
	GlobalReactions []DiscussionReaction        `json:"discussionReaction"`
}

type DiscussionErrorResponse struct {
	TimeStamp timeutil.TimeStamp `json:"timestamp"`
	Status    int                `json:"status"`
	Error     string             `json:"error"`
	message   string
}

func handleResponse(resp *resty.Response) error {
	if resp.IsError() {
		var errResp DiscussionErrorResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			log.Error("Failed to parse error response: %v", err)
			return fmt.Errorf("unexpected error: %s", resp.Status())
		}
		log.Error("API Error %d: %s", errResp.Status, errResp.message)
		return fmt.Errorf("api error %d: %s", errResp.Status, errResp.message)
	}
	return nil
}

func PostDiscussion(request *PostDiscussionRequest) (int, error) {
	resp, err := client.Request().
		SetBody(request).
		Post("/discussion")
	if err != nil {
		return -1, fmt.Errorf("failed to make POST /discussion request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return -1, err
	}

	var result int
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return -1, fmt.Errorf("failed to parse response body: %w", err)
	}

	return result, nil
}

func GetDiscussion(repoId int64) (*DiscussionResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d request: %w", repoId, err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var result DiscussionResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionCount(repoId int64) (*DiscussionCountResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d/count", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/count request: %w", repoId, err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var result DiscussionCountResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionList(repoId int64, isClosed bool, page int) (*DiscussionListResponse, error) {
	isClosedAsString := strconv.FormatBool(isClosed)
	pageAsString := strconv.Itoa(page)
	resp, err := client.Request().
		SetQueryParam("isClosed", isClosedAsString).
		SetQueryParam("page", pageAsString).
		Get(fmt.Sprintf("/discussion/%d/list", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/list request: %w", repoId, err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var result DiscussionListResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionContent(discussionId int64) (*DiscussionContentResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d/contents", discussionId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/contents request: %w", discussionId, err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var result DiscussionContentResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionComment(discussionCommentId int64) (*DiscussionCommentResponse, error) {
	resp, err := client.Request().
		SetQueryParam("id", strconv.FormatInt(discussionCommentId, 10)).
		Get("/discussion/comment")
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/comment request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var result DiscussionCommentResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func HandleDiscussionAvailable() (*resty.Response, error) {
	resp, err := client.Request().
		Post("/discussion/available")
	if err != nil {
		return nil, fmt.Errorf("failed to make POST /discussion/available request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func PostComment(request *PostCommentRequest) (*int64, error) {
	resp, err := client.Request().
		SetBody(request).
		Post("/discussion/comment")
	if err != nil {
		return nil, fmt.Errorf("failed to make POST /discussion/comment request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	bodyStr := string(resp.Body())
	id, err := strconv.ParseInt(bodyStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &id, nil
}

func ModifyDiscussion(request *ModifyDiscussionRequest) (*resty.Response, error) {
	resp, err := client.Request().
		SetBody(request).
		Put("/discussion")
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT /discussion request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func DeleteDiscussionComment(discussionCommentId int64, posterId int64) error {
	request := &DeleteDiscussionCommentRequest{
		DiscussionCommentId: discussionCommentId,
		PosterId:            posterId,
	}
	resp, err := client.Request().
		SetBody(request).
		Delete("/discussion/comment")
	if err != nil {

		return fmt.Errorf("failed to make DELETE /discussion/comment request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return err
	}

	return nil
}

func ModifyDiscussionComment(request *ModifyDiscussionCommentRequest) error {
	resp, err := client.Request().
		SetBody(request).
		Put("/discussion/comment")
	if err != nil {
		return fmt.Errorf("failed to make PUT /discussion/comment request: %w", err)
	}

	if err := handleResponse(resp); err != nil {
		return err
	}

	return nil
}

/**
 * discussion methods
 * the `discussion` struct could be moved to a separate file later
 */
func (d *Discussion) IsExpired() bool {
	return d.DeadlineUnix < timeutil.TimeStamp(time.Now().Unix())
}

func (d *Discussion) GetLastEventTimestamp() timeutil.TimeStamp {
	if d.IsClosed {
		return d.ClosedUnix
	}
	return d.CreatedUnix
}

func (d *Discussion) GetLastEventLabel() string {
	if d.IsClosed {
		return "repo.discussion.closed_by"
	}
	return "repo.discussion.opened_by"
}

func (d *Discussion) GetLastEventLabelFake() string {
	if d.IsClosed {
		return "repo.discussion.closed_by_fake"
	}
	return "repo.discussion.opened_by_fake"
}

func (d *Discussion) LoadPoster(ctx *context.Context) (err error) {
	if d.Poster == nil && d.PosterId != 0 {
		d.Poster, err = user_model.GetPossibleUserByID(ctx, d.PosterId)
		if err != nil {
			d.PosterId = user_model.GhostUserID
			d.Poster = user_model.NewGhostUser()
			if !user_model.IsErrUserNotExist(err) {
				return fmt.Errorf("getUserById.(poster) [%d]: %w", d.PosterId, err)
			}
			return nil
		}
	}
	return err
}

func (d *Discussion) LoadRepo(ctx *context.Context) (err error) {
	d.Repo = ctx.Repo.Repository
	return nil
}

func GetDiscussionContents(discussionId int64) (*DiscussionContentResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d/contents", discussionId))
	if err != nil {
		return nil, err
	}
	result := &DiscussionContentResponse{}
	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return nil, err
	}

	if err := handleResponse(resp); err != nil {
		return nil, fmt.Errorf("failed to GetDiscussionContents %w", err)
	}

	return result, nil
}

func SetDiscussionClosedState(discussionId int64, isClosed bool) error {
	resp, err := client.Request().Patch(fmt.Sprintf("discussion/state/%d?isClosed=%t", discussionId, isClosed))
	if err != nil {
		return err
	}

	if err := handleResponse(resp); err != nil {
		return fmt.Errorf("failed to set review state, got %d", resp.StatusCode())
	}

	return nil
}


func (dr DiscussionResponse) IsPoster(id int64) bool {
	return dr.PosterId == id
}
