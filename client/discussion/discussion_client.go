package discussion

import (
	"database/sql"
	"fmt"

	"code.gitea.io/gitea/client"
	"github.com/go-resty/resty/v2"
)

type DiscussionCode struct {
	Id        sql.NullInt64 `json:"id"`
	FilePath  string        `json:"filePath"`
	StartLine int           `json:"startLine"`
	EndLine   int           `json:"endLine"`
}

type PostDiscussionRequest struct {
	RepoId   int64            `json:"repoId"`
	PosterId int64            `json:"posterId"`
	Name     string           `json:"name"`
	Content  string           `json:"content"`
	Codes    []DiscussionCode `json:"codes"`
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

func PostDiscussion(request PostDiscussionRequest) (*resty.Response, error) {
	return client.Request().SetBody(request).Post("/discussion")
}

func GetDiscussionCount(repoId int64, isClosed bool) (*resty.Response, error) {
	var isClosedAsInt = map[bool]int{false: 0, true: 1}[isClosed]
	return client.Request().
		SetQueryParam("isClosed", string(isClosedAsInt)).
		Get(fmt.Sprintf("/discussion/%d/count", repoId))
}

func HandleDiscussionAvailable() (*resty.Response, error) {
	return client.Request().Post("/discussion/available")
}

func GetDiscussionContents(discussionId int64) (*resty.Response, error) {
	return client.Request().Get(fmt.Sprintf("/discussion/%d/codes", discussionId))
}

func PostComment(request PostCommentRequest) (*resty.Response, error) {
	return client.Request().SetBody(request).Post("/discussion/comment")
}

func ModifyDiscussion(request ModifyDiscussionRequest) (*resty.Response, error) {
	return client.Request().SetBody(request).Put("/discussion")
}
