package recorder

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
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

/**
Request{
	ctx:        ctx,
	Method:     method,
	URL:        u,
	Proto:      "HTTP/1.1",
	ProtoMajor: 1,
	ProtoMinor: 1,
	Header:     make(Header),
	Body:       rc,
	Host:       u.Host,
}
*/

func (this *Request) FromHTTPRequest(httpReq *http.Request) *Request {
	body := &bytes.Buffer{}
	httpReq.Body = ioutil.NopCloser(io.TeeReader(httpReq.Body, body))

	this.Proto = httpReq.Proto
	this.Method = httpReq.Method
	this.URL = httpReq.URL.String()
	this.Path = httpReq.URL.Path
	this.Header = httpReq.Header
	this.Body = body.Bytes()

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
