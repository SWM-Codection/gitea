package git

import (
	"fmt"

	"code.gitea.io/gitea/client"
	"github.com/go-resty/resty/v2"
)

func ListGitFiles(userName string, repoName string) (*resty.Response, error) {
	return client.Request().
		Get(fmt.Sprintf("%s/%s/discussions", userName, repoName))
}

func GetFileContent(userName string, repoName string, filePath string) (*resty.Response, error) {
	return client.Request().
		SetQueryParam("filepath", filePath).
		Get(fmt.Sprintf("%s/%s/discussions/contents", userName, repoName))
}
