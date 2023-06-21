package tunnel

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/ecromaneli-golang/http/webserver"
)

func (req *Request) ToString() string {
	return "[" + req.Method + "] " + req.Path
}

func (req *Request) FromHTTPRequest(wsReq *webserver.Request) *Request {
	req.Proto = wsReq.Raw.Proto
	req.Method = wsReq.Raw.Method
	req.Url = wsReq.Raw.URL.String()
	req.Path = wsReq.Raw.URL.Path
	req.Body = wsReq.Body()
	req.setHttpHeader(wsReq.Raw.Header)
	req.setUserIP(wsReq)

	return req
}

func (req *Request) ToHTTPRequest() *http.Request {
	var data io.Reader = nil

	if req.Body != nil {
		data = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.Url, data)

	if err != nil {
		panic(err)
	}

	httpReq.Proto = req.Proto
	httpReq.Header = req.GetHttpHeader()

	return httpReq
}

func (req *Request) setUserIP(wsReq *webserver.Request) {
	ipAddress := wsReq.Raw.Header.Get("X-Real-Ip")

	if ipAddress == "" {
		ipAddress = wsReq.Raw.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = wsReq.Raw.RemoteAddr
	}

	req.RemoteAddr = ipAddress
}

func (req *Request) GetHttpHeader() http.Header {
	httpHeader := http.Header{}

	for key, value := range req.Header {
		httpHeader[key] = strings.Split(value, ",")
	}

	return httpHeader
}

func (req *Request) setHttpHeader(header http.Header) {
	req.Header = make(map[string]string)

	for key, values := range header {
		req.Header[key] = strings.Join(values, ",")
	}
}

func (req *Request) Clone() *Request {
	clone := &Request{
		Proto:      req.Proto,
		Method:     req.Method,
		Url:        req.Url,
		Path:       req.Path,
		Header:     make(map[string]string),
		Body:       req.Body,
		RemoteAddr: req.RemoteAddr,
	}

	for key, value := range req.Header {
		clone.Header[key] = value
	}

	return clone
}
