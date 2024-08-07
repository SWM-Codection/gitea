package discussion

import (
	"database/sql"
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
	CodeId       int64            `json:"codeId"`
	PosterId     int64            `json:"posterId"`
	Scope        CommentScopeEnum `json:"scope"`
	StartLine    sql.NullInt32    `json:"startLine"`
	EndLine      sql.NullInt32    `json:"endLine"`
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
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Content    string  `json:"content"`
	RepoId     int64   `json:"repoId"`
	PosterId   int64   `json:"posterId"`
	CommitHash string  `json:"commitHash"`
	IsClosed   bool    `json:"isClosed"`
	Deadline   string  `json:"deadline"`
	Assignees  []int64 `json:"assignees"`
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

type DiscussionCountResponse struct {
	OpenCount  int `json:"openCount"`
	CloseCount int `json:"closeCount"`
}

func PostDiscussion(request *PostDiscussionRequest) (int, error) {
	log.Info("PostDiscussion request : %v", request)
	resp, err := client.Request().SetBody(request).Post("/discussion")
	if err != nil {
		return -1, err
	}
	log.Info("resp.string() : %v", resp.String())
	result, err := strconv.Atoi(resp.String())
	if err != nil {
		return -1, err
	}
	return result, err
}

func GetDiscussion(repoId int64) (*DiscussionResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d", repoId))
	if err != nil {
		return nil, err
	}
	var result = &DiscussionResponse{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetDiscussionCount(repoId int64) (*DiscussionCountResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d/count", repoId))
	if err != nil {
		return nil, err
	}
	result := &DiscussionCountResponse{}
	if err = json.Unmarshal(resp.Body(), result); err != nil {
		return nil, err
	}
	return result, nil
}

func GetDiscussionList(repoId int64, isClosed bool, page int) (*DiscussionListResponse, error) {
	isClosedAsString := strconv.FormatBool(isClosed)
	pageAsString := strconv.Itoa(page)
	resp, err := client.Request().
		SetQueryParam("isClosed", isClosedAsString).
		SetQueryParam("page", pageAsString).
		Get(fmt.Sprintf("/discussion/%d/list", repoId))
	if err != nil {
		return nil, err
	}
	result := &DiscussionListResponse{}
	json.Unmarshal(resp.Body(), result)
	return result, nil
}

func HandleDiscussionAvailable() (*resty.Response, error) {
	return client.Request().Post("/discussion/available")
}

func GetDiscussionContents(discussionId int64) (*resty.Response, error) {
	return client.Request().Get(fmt.Sprintf("/discussion/%d/codes", discussionId))
}

func PostComment(request *PostCommentRequest) (*resty.Response, error) {
	return client.Request().SetBody(request).Post("/discussion/comment")
}

func ModifyDiscussion(request *ModifyDiscussionRequest) (*resty.Response, error) {
	return client.Request().SetBody(request).Put("/discussion")
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
