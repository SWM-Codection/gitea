package client

import (
	"sync"

	"code.gitea.io/gitea/modules/setting"
	"github.com/go-resty/resty/v2"
)

var (
	client *resty.Client
	once   sync.Once
)

func Request() *resty.Request {
	once.Do(func() {
		client = resty.New().SetBaseURL(setting.DiscussionServer.Url)
	})
	return client.R()
}
