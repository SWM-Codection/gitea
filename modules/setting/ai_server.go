package setting

import "fmt"

var AiServer = struct {
	Host string
	Port int
	Url  string
}{
	Host: "localhost",
	Port: 8000,
	Url:  "http://localhost:8000",
}

func loadAiServerFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("ai_server")
	AiServer.Host = sec.Key("host").MustString("localhost")
	AiServer.Port = sec.Key("PORT").MustInt(8000)
	AiServer.Url = fmt.Sprintf("http://%s:%d", AiServer.Host, AiServer.Port)
}
