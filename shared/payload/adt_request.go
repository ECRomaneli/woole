package payload

import (
	"bytes"
	"io"
	"net/http"

	"github.com/ecromaneli-golang/http/webserver"
)

type Request struct {
	Proto  string
	Method string
	URL    string
	Path   string
	Header http.Header
	Body   []byte
}

func (this *Request) ToString() string {
	return "[" + this.Method + "] " + this.Path
}

func (this *Request) FromHTTPRequest(httpReq *webserver.Request) *Request {
	this.Proto = httpReq.Raw.Proto
	this.Method = httpReq.Raw.Method
	this.URL = httpReq.Raw.URL.String()
	this.Path = httpReq.Raw.URL.Path
	this.Header = httpReq.Raw.Header
	this.Body = httpReq.Body()

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
