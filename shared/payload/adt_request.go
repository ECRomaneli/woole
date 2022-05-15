package payload

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ecromaneli-golang/http/webserver"
)

type Request struct {
	Proto      string      `json:"proto"`
	Method     string      `json:"method"`
	URL        string      `json:"url"`
	Path       string      `json:"path"`
	Header     http.Header `json:"header"`
	Body       []byte      `json:"body"`
	RemoteAddr string      `json:"remoteAddr"`
}

func (this *Request) ToString() string {
	return "[" + this.Method + "] " + this.Path
}

func (this *Request) FromHTTPRequest(req *webserver.Request) *Request {
	this.Proto = req.Raw.Proto
	this.Method = req.Raw.Method
	this.URL = req.Raw.URL.String()
	this.Path = req.Raw.URL.Path
	this.Header = req.Raw.Header
	this.Body = req.Body()
	this.setUserIP(req)

	return this
}

func (this *Request) ToHTTPRequest() *http.Request {
	var data io.Reader = nil

	if this.Body != nil {
		data = bytes.NewReader(this.Body)
	}

	httpReq, err := http.NewRequest(this.Method, this.URL, data)

	if err != nil {
		panic(err)
	}

	httpReq.Proto = this.Proto
	httpReq.Header = this.Header

	return httpReq
}

func (this *Request) setUserIP(req *webserver.Request) {
	ipAddress := req.Raw.Header.Get("X-Real-Ip")

	if ipAddress == "" {
		ipAddress = req.Raw.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = req.Raw.RemoteAddr
	}

	this.RemoteAddr = ipAddress
}
