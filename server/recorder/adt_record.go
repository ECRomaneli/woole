package recorder

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"woole-server/util/signal"
)

type Record struct {
	Id         string
	Request    *Request
	Response   *Response
	Elapsed    time.Duration
	OnResponse *signal.Signal
}

type Records struct {
	sync.RWMutex
	clients map[string]*Client
}

func NewRecords() *Records {
	return &Records{clients: make(map[string]*Client)}
}

func NewRecord(req *Request) *Record {
	return &Record{
		Request:    req,
		OnResponse: signal.New(),
	}
}

func (this *Records) RegisterClient(clientId string) *Client {
	if this.ClientExists(clientId) {
		this.RemoveClient(clientId)
	}

	this.clients[clientId] = NewClient(clientId)
	return this.clients[clientId]
}

func (this *Records) ClientExists(clientId string) bool {
	return this.clients[clientId] != nil
}

func (this *Records) ClientIsLocked(clientId string) bool {
	return this.clients[clientId] != nil && this.clients[clientId].IsLocked()
}

func (this *Records) Add(clientId string, rec *Record) (id string) {
	this.RLock()
	client := this.clients[clientId]
	this.RUnlock()

	return client.Add(rec)
}

func (this *Records) Remove(clientId, recordId string) *Record {
	this.RLock()
	client := this.clients[clientId]
	this.RUnlock()

	if client == nil {
		return nil
	}

	return client.Remove(recordId)
}

func (this *Records) RemoveClient(clientId string) {
	this.Lock()
	defer this.Unlock()

	close(this.clients[clientId].Tunnel)
	this.clients[clientId] = nil
}

func (this *Records) Get(clientId string, bearer string) (*Client, error) {
	this.RLock()
	defer this.RUnlock()

	client := this.clients[clientId]

	if client.Authorize(bearer) {
		return client, nil
	}

	return nil, errors.New("Authentication failed for client '" + clientId + "'")
}

func (this *Record) ToString(maxPathLength int) string {
	path := []byte(this.Request.Path)

	if len(path) > maxPathLength {
		path = append([]byte("..."), path[len(path)-maxPathLength:]...)
	}

	method := "[" + this.Request.Method + "]"

	strPathLength := strconv.Itoa(maxPathLength + 3)
	str := fmt.Sprintf("%8s %"+strPathLength+"s", method, string(path))

	if this.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", this.Response.Code, this.Elapsed)
}
