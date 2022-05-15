package payload

import (
	"net/http"
	"net/http/httptest"
)

type Response struct {
	Proto  string      `json:"proto"`
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Header http.Header `json:"header"`
	Body   []byte      `json:"body"`
}

func (this *Response) FromResponseRecorder(httpRes *httptest.ResponseRecorder) *Response {

	this.Proto = httpRes.Result().Proto
	this.Status = httpRes.Result().Status
	this.Code = httpRes.Code
	this.Header = httpRes.Header()
	this.Body = httpRes.Body.Bytes()

	return this
}
