package payload

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func (res *Response) FromResponseRecorder(httpRes *httptest.ResponseRecorder, elapsed int64) *Response {

	res.Proto = httpRes.Result().Proto
	res.Status = httpRes.Result().Status
	res.Code = int32(httpRes.Code)
	res.Body = httpRes.Body.Bytes()
	res.Elapsed = elapsed
	res.SetHttpHeader(httpRes.Header())

	return res
}

func (res *Response) GetHttpHeader() http.Header {
	httpHeader := http.Header{}

	for key, value := range res.Header {
		httpHeader[key] = strings.Split(value, ",")
	}

	return httpHeader
}

func (res *Response) SetHttpHeader(header http.Header) {
	res.Header = make(map[string]string)

	for key, values := range header {
		res.Header[key] = strings.Join(values, ",")
	}
}

func (res *Response) Clone() *Response {
	clone := &Response{
		Proto:   res.Proto,
		Status:  res.Status,
		Code:    res.Code,
		Header:  make(map[string]string),
		Body:    res.Body,
		Elapsed: res.Elapsed,
	}

	for key, value := range res.Header {
		clone.Header[key] = value
	}

	return clone
}
