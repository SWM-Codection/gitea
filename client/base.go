package client

import (
	"sync"

	"code.gitea.io/gitea/modules/setting"
	"github.com/go-resty/resty/v2"
)

var client *resty.Client = nil
var clientMtx = &sync.Mutex{}

func Request() *resty.Request {
	clientMtx.Lock()
	if client == nil {
		client = resty.New().SetBaseURL(setting.DiscussionServer.Url)
	}
	defer clientMtx.Unlock()
	return client.R()
}
