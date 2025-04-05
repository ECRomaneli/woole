package tunnel

import "woole/internal/pkg/url"

func (auth *Session) HttpUrl() string {
	var port string

	if url.GetDefaultPortStr("http") != auth.GetHttpPort() {
		port = ":" + auth.GetHttpPort()
	}

	return "http://" + auth.GetHostname() + port
}

func (auth *Session) HttpsUrl() string {
	if auth.HttpsPort == "" {
		return ""
	}

	var port string

	if url.GetDefaultPortStr("https") != auth.GetHttpsPort() {
		port = ":" + auth.GetHttpsPort()
	}

	return "https://" + auth.GetHostname() + port
}
