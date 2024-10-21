package discussion

import (
	"encoding/json"
	"fmt"
	"strconv"

	"code.gitea.io/gitea/client"
	"code.gitea.io/gitea/client/discussion/model"
	"code.gitea.io/gitea/modules/log"
	"github.com/go-resty/resty/v2"
)

func validateResponse(resp *resty.Response) error {
	if resp.IsError() {
		var errResp model.DiscussionErrorResponse
		if err := json.Unmarshal(resp.Body(), &errResp); err != nil {
			log.Error("Failed to parse error response: %v", err)
			return fmt.Errorf("unexpected error: %s", resp.Status())
		}
		log.Error("API Error %d: %s", errResp.Status, errResp.Message)
		return fmt.Errorf("api error %d: %s", errResp.Status, errResp.Message)
	}
	return nil
}

func PostDiscussion(request *model.PostDiscussionRequest) (int, error) {
	resp, err := client.Request().
		SetBody(request).
		Post("/discussion")
	if err != nil {
		return -1, fmt.Errorf("failed to make POST /discussion request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return -1, err
	}

	var result int
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return -1, fmt.Errorf("failed to parse response body: %w", err)
	}

	return result, nil
}

func GetDiscussion(repoId int64) (*model.DiscussionResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d request: %w", repoId, err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	var result model.DiscussionResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionCount(repoId int64) (*model.DiscussionCountResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d/count", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/count request: %w", repoId, err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	var result model.DiscussionCountResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionList(repoId int64, isClosed bool, page int, sort string) (*model.DiscussionListResponse, error) {
	isClosedAsString := strconv.FormatBool(isClosed)
	pageAsString := strconv.Itoa(page)
	resp, err := client.Request().
		SetQueryParam("isClosed", isClosedAsString).
		SetQueryParam("page", pageAsString).
		SetQueryParam("sort", sort).
		Get(fmt.Sprintf("/discussion/%d/list", repoId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/list request: %w", repoId, err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	var result model.DiscussionListResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionContent(discussionId int64) (*model.DiscussionContentResponse, error) {
	resp, err := client.Request().
		Get(fmt.Sprintf("/discussion/%d/contents", discussionId))
	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/%d/contents request: %w", discussionId, err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	var result model.DiscussionContentResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionComment(discussionCommentId int64) (*model.DiscussionCommentResponse, error) {

	resp, err := client.Request().
		SetQueryParam("id", strconv.FormatInt(discussionCommentId, 10)).
		Get("/discussion/comment")

	if err != nil {
		return nil, fmt.Errorf("failed to make GET /discussion/comment request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	var result model.DiscussionCommentResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &result, nil
}

func GetDiscussionCommentsByCodeId(codeId int64) ([]*model.DiscussionCommentResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/comments/%d", codeId))

	result := make([]*model.DiscussionCommentResponse, 0)

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)

	}

	return result, nil
}

func HandleDiscussionAvailable() (*resty.Response, error) {
	resp, err := client.Request().
		Post("/discussion/available")
	if err != nil {
		return nil, fmt.Errorf("failed to make POST /discussion/available request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func PostComment(request *model.PostCommentRequest) (*int64, error) {
	resp, err := client.Request().
		SetBody(request).
		Post("/discussion/comment")
	if err != nil {
		return nil, fmt.Errorf("failed to make POST /discussion/comment request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyStr := string(resp.Body())
	id, err := strconv.ParseInt(bodyStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return &id, nil
}

func ModifyDiscussion(request *model.ModifyDiscussionRequest) (*resty.Response, error) {
	resp, err := client.Request().
		SetBody(request).
		Put("/discussion")
	if err != nil {
		return nil, fmt.Errorf("failed to make PUT /discussion request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func DeleteDiscussionComment(discussionCommentId int64, posterId int64) error {
	request := &model.DeleteDiscussionCommentRequest{
		DiscussionCommentId: discussionCommentId,
		PosterId:            posterId,
	}
	resp, err := client.Request().
		SetBody(request).
		Delete("/discussion/comment")
	if err != nil {

		return fmt.Errorf("failed to make DELETE /discussion/comment request: %w", err)
	}

	if err := validateResponse(resp); err != nil {
		return err
	}

	return nil
}

func GetDiscussionContents(discussionId int64) (*model.DiscussionContentResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d/contents", discussionId))
	if err != nil {
		return nil, err
	}
	result := &model.DiscussionContentResponse{}

	if err := json.Unmarshal(resp.Body(), result); err != nil {
		return nil, err
	}

	if err := validateResponse(resp); err != nil {
		return nil, fmt.Errorf("failed to GetDiscussionContents %w", err)
	}

	return result, nil
}

func SetDiscussionClosedState(discussionId int64, isClosed bool) error {
	resp, err := client.Request().Patch(fmt.Sprintf("discussion/state/%d?isClosed=%t", discussionId, isClosed))
	if err != nil {
		return err
	}

	if err := validateResponse(resp); err != nil {
		return fmt.Errorf("failed to set review state, got %d", resp.StatusCode())
	}

	return nil
}

func SetDiscussionDeadline(discussionId int64, deadline int64) error {
	resp, err := client.Request().SetQueryParam("deadline", strconv.FormatInt(deadline, 10)).Patch(fmt.Sprintf("/discussion/deadline/%d", discussionId))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to set deadline, got %d", resp.StatusCode())
	}

	return nil
}

func UpdateDiscussionAssignee(request *model.UpdateAssigneeRequest) error {
	resp, err := client.Request().SetBody(request).Put("/discussion/assignees")
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to update assignee, got %d", resp.StatusCode())
	}
	return nil
}

func ModifyDiscussionComment(request *model.ModifyDiscussionCommentRequest) error {
	resp, err := client.Request().SetBody(request).Put("/discussion/comment")
	if err != nil {
		return err
	}
	if err := validateResponse(resp); err != nil {
		return err
	}
	return nil
}

func ClearDiscussionAssignee(discussionId int64) error {
	resp, err := client.Request().Delete(fmt.Sprintf("/discussion/assignees/%d", discussionId))
	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("failed to clear assignee, got %d", resp.StatusCode())
	}

	return nil
}

func GiveReaction(request model.DiscussionReactionRequest) (int64, error) {
	var result int64 = -1
	resp, err := client.Request().SetBody(request).Post("/discussion/reaction")
	if err != nil {
		return result, err
	}
	if err := validateResponse(resp); err != nil {
		return result, err
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return result, err
	}
	return result, err
}

func RemoveReaction(request model.DiscussionReactionRequest) error {
	resp, err := client.Request().SetBody(request).Delete("/discussion/reaction")
	if err != nil {
		return err
	}
	if err := validateResponse(resp); err != nil {
		return err
	}
	return err
}

func GetDiscussionCommentReaction(commentId int64) (*model.ReactionList, error) {
	resp, err := client.Request().SetQueryParam("commentId", strconv.FormatInt(commentId, 10)).
		Get("/discussion/reaction")
	if err != nil {
		log.Error("Failed to make GET /discussion/reaction request: %v", err)
		return nil, fmt.Errorf("failed to make GET /discussion/reaction request: %w", err)
	}
	var result model.ReactionList
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	return &result, err
}

func IsNewPinAllowed(repoId int64) (bool, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d/max-pin", repoId))
	if err != nil {
		return false, err
	}

	if err := validateResponse(resp); err != nil {
		return false, err
	}
	bodyStr := string(resp.Body())

	isAllowed, err := strconv.ParseBool(string(bodyStr))
	if err != nil {
		return false, fmt.Errorf("failed to parse response body: %w", err)
	}

	return isAllowed, nil
}

func ConvertDiscussionPinStatus(discussionId int64) error {
	resp, err := client.Request().Post(fmt.Sprintf("discussion/%d/pin", discussionId))
	if err != nil {
		return err
	}

	if err := validateResponse(resp); err != nil {
		return fmt.Errorf("failed to convert discussion pin state, got %d", resp.StatusCode())
	}

	return nil
}

func UnpinDiscussion(discussionId int64) error {
	resp, err := client.Request().Delete(fmt.Sprintf("discussion/%d/unpin", discussionId))
	if err != nil {
		return err
	}

	if err := validateResponse(resp); err != nil {
		return fmt.Errorf("failed to unpin discussion, got %d", resp.StatusCode())
	}

	return nil
}

func GetPinnedDiscussions(repoId int64) (*model.DiscussionListResponse, error) {
	resp, err := client.Request().Get(fmt.Sprintf("/discussion/%d/pin", repoId))
	if err != nil {
		return nil, err
	}

	if err := validateResponse(resp); err != nil {
		return nil, err
	}
	var result model.DiscussionListResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}
	return &result, nil
}

func MoveDiscussionPin(request *model.MoveDiscussionPinRequest) error {
	resp, err := client.Request().
		SetBody(request).
		Post("discussion/move-pin")
	if err != nil {
		return fmt.Errorf("failed to make POST /discussion request: %w", err)
	}
	if err := validateResponse(resp); err != nil {
		return err
	}

	return nil
}
