package payload

import (
	"woole/shared/util"
)

func (auth *Session) HttpUrl() string {
	var port string

	if util.GetDefaultPortStr("http") != auth.GetHttpPort() {
		port = ":" + auth.GetHttpPort()
	}

	return "http://" + auth.GetHostname() + port
}

func (auth *Session) HttpsUrl() string {
	if auth.HttpsPort == "" {
		return ""
	}

	var port string

	if util.GetDefaultPortStr("https") != auth.GetHttpsPort() {
		port = ":" + auth.GetHttpsPort()
	}

	return "https://" + auth.GetHostname() + port
}
