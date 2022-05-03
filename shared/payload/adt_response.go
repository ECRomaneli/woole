package payload

import (
	"net/http"
	"net/http/httptest"
)

type Response struct {
	Proto  string
	Status string
	Code   int
	Header http.Header
	Body   []byte
}

func (this *Response) FromResponseRecorder(httpRes *httptest.ResponseRecorder) *Response {

	this.Proto = httpRes.Result().Proto
	this.Status = httpRes.Result().Status
	this.Code = httpRes.Code
	this.Header = httpRes.Header()
	this.Body = httpRes.Body.Bytes()

	return this
}
