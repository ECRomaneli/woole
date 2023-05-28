package adt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

const pathParamPrefix = "/woole?param="

type RedirectAction int8

const (
	CONTINUE        RedirectAction = 0
	CHANGE_URL_HOST RedirectAction = 1
	// TODO: CHANGE_CLIENT_HOST RedirectAction = 2
)

type Redirect struct {
	RecordId string `json:"recordId,omitempty"`
	/*Url      string         `json:"url,omitempty"`*/
	Action RedirectAction `json:"action,omitempty"`
}

type CustomPathParam struct {
	Redirect Redirect `json:"redirect,omitempty"`
}

func (customPathParam *CustomPathParam) Serialize() string {
	jsonParam, err := json.Marshal(*customPathParam)
	if err != nil {
		panic(err)
	}
	return pathParamPrefix + string(base64.StdEncoding.EncodeToString(jsonParam))
}

func DeserializeCustomPathParam(path string) (param *CustomPathParam, ok bool) {
	b64Param, found := strings.CutPrefix(path, pathParamPrefix)

	if !found {
		return nil, false
	}

	jsonParam, err := base64.StdEncoding.DecodeString(b64Param)
	if err != nil {
		panic(err)
	}

	customPathParam := &CustomPathParam{}

	err = json.Unmarshal(jsonParam, customPathParam)
	if err != nil {
		panic(err)
	}
	return customPathParam, true
}
