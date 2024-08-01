package client

import (
	"code.gitea.io/gitea/modules/setting"
	"github.com/go-resty/resty/v2"
)

var client *resty.Client = nil

func Request() *resty.Request {
	if client == nil {
		client = resty.New().SetBaseURL(setting.DiscussionServer.Url)
	}
	return client.R()
}
