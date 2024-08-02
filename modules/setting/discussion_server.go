package setting

import "fmt"

var DiscussionServer = struct {
	Host string
	Port int
	Url  string
}{
	Host: "localhost",
	Port: 8081,
	Url:  "http://localhost:8081",
}

func loadDiscussionServerFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("discussion_server")
	DiscussionServer.Host = sec.Key("HOST").MustString("localhost")
	DiscussionServer.Port = sec.Key("PORT").MustInt(8081)
	DiscussionServer.Url = fmt.Sprintf("http://%s:%d", DiscussionServer.Host, DiscussionServer.Port)
}
