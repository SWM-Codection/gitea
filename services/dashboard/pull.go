package dashboard

// import (
// 	"code.gitea.io/gitea/modules/cache"
// 	"code.gitea.io/gitea/modules/graceful"
// 	"code.gitea.io/gitea/services/context"
// 	"errors"
// 	"fmt"
// 	"golang.org/x/exp/maps"
// 	"sync"
// 	"time"
// )

// const (
// 	FeedBackTimesKey                   = "GetFeedBackTimeStatus/%d"
// 	dateLayout                         = time.DateOnly
// 	contributorStatsCacheTimeout int64 = 60 * 10
// )

// var (
// 	ErrAwaitGeneration  = errors.New("generation took longer than ")
// 	awaitGenerationTime = time.Second * 5
// 	generateLock        = sync.Map{}
// )

// func findLastSundayBeforeDate(dateStr string) (string, error) {
// 	date, err := time.Parse(dateLayout, dateStr)
// 	if err != nil {
// 		return "", err
// 	}

// 	weekday := date.Weekday()
// 	daysToSubtract := int(weekday) - int(time.Sunday)
// 	if daysToSubtract < 0 {
// 		daysToSubtract += 7
// 	}

// 	lastSunday := date.AddDate(0, 0, -daysToSubtract)
// 	return lastSunday.Format(dateLayout), nil
// }

// type FeedBackTimeStatus struct {
// 	Week int64 `json:"week"`
// 	Diff int64 `json:"feedback_time"`
// }

// type DashboardService interface {
// 	GetFeedBackTimeStatus(ctx *context.Context, OrgID int64, cache cache.StringCache) ([]*FeedBackTimeStatus, error)
// 	getFeedBackTimeStatus(genDone chan struct{}, cache cache.StringCache, cacheKey string, OrgID int64)
// }

// type DashboardServiceImpl struct{}

// var _ DashboardService = &DashboardServiceImpl{}

// func (is *DashboardServiceImpl) getFeedBackTimeStatus(genDone chan struct{}, cache cache.StringCache, cacheKey string, OrgID int64) {

// 	ctx := graceful.GetManager().HammerContext()

// 	repositoryIDs, err := DpullDbAdapter.GetOrgRepositoryOrgIDS(&ctx, OrgID)

// 	if err != nil {
// 		return
// 	}

// 	firstCreatedUnixes, err := DpullDbAdapter.GetFirstReviewCreatedUnixesByRepoIDs(&ctx, repositoryIDs)

// 	if err != nil {
// 		return
// 	}

// 	feedBackTimeStatus := make(map[int64]*FeedBackTimeStatus)

// 	for _, firstCreatedUnix := range firstCreatedUnixes {

// 		diff := firstCreatedUnix.CalDiff()

// 		date, _ := findLastSundayBeforeDate(firstCreatedUnix.IssueCreated.FormatDate())

// 		val, _ := time.Parse(dateLayout, date)

// 		week := val.UnixMilli()

// 		if feedBackTimeStatus[week] == nil {
// 			feedBackTimeStatus[week] = &FeedBackTimeStatus{
// 				Week: week,
// 				Diff: 0,
// 			}

// 		}

// 		feedBackTimeStatus[week].Diff += diff.AsTime().UnixMilli()

// 	}

// 	if err = cache.PutJSON(cacheKey, maps.Values(feedBackTimeStatus), contributorStatsCacheTimeout); err != nil {
// 		return
// 	}

// 	generateLock.Delete(cacheKey)
// 	if genDone != nil {
// 		genDone <- struct{}{}
// 	}
// }

// func (is *DashboardServiceImpl) GetFeedBackTimeStatus(ctx *context.Context, OrgID int64, cache cache.StringCache) ([]*FeedBackTimeStatus, error) {

// 	cacheKey := fmt.Sprintf(FeedBackTimesKey, OrgID)

// 	if !cache.IsExist(cacheKey) {
// 		genReady := make(chan struct{})

// 		_, run := generateLock.Load(cacheKey)

// 		if run {
// 			return nil, ErrAwaitGeneration
// 		}

// 		generateLock.Store(cacheKey, struct{}{})

// 		go is.getFeedBackTimeStatus(genReady, cache, cacheKey, OrgID)

// 		select {
// 		case <-time.After(awaitGenerationTime):
// 			return nil, ErrAwaitGeneration
// 		case <-genReady:

// 			break
// 		}

// 	}

// 	var res []*FeedBackTimeStatus

// 	_, err := cache.GetJSON(cacheKey, &res)

// 	if err != nil {
// 		return nil, fmt.Errorf("GetFeedBackTimeStatus: %v", err.ToError())
// 	}

// 	return res, nil
// }
