package payload

import (
	"net/http"
	"net/http/httptest"
)

func (res *Response) FromResponseRecorder(httpRes *httptest.ResponseRecorder) *Response {

	res.Proto = httpRes.Result().Proto
	res.Status = httpRes.Result().Status
	res.Code = int32(httpRes.Code)
	res.Body = httpRes.Body.Bytes()
	res.setHttpHeader(httpRes.Header())

	return res
}

func (res *Response) GetHttpHeader() http.Header {
	httpHeader := http.Header{}

	for key, stringList := range res.Header {
		if stringList == nil {
			httpHeader[key] = []string{}
			continue
		}
		// if null > res.Header[key] = &StringList{}
		httpHeader[key] = stringList.Val
	}

	return httpHeader
}

func (res *Response) setHttpHeader(header http.Header) {
	res.Header = map[string]*StringList{}

	for key, values := range header {
		res.Header[key] = &StringList{Val: values}
	}
}
