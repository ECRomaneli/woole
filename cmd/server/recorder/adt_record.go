package recorder

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"woole/shared/payload"
	"woole/shared/util/signal"
)

type Record struct {
	Id         string
	Request    *payload.Request
	Response   *payload.Response
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

func NewRecord(req *payload.Request) *Record {
	return &Record{
		Request:    req,
		OnResponse: signal.New(),
	}
}

func (recs *Records) RegisterClient(clientId string) *Client {
	if recs.ClientExists(clientId) {
		recs.RemoveClient(clientId)
	}

	recs.clients[clientId] = NewClient(clientId)
	return recs.clients[clientId]
}

func (recs *Records) ClientExists(clientId string) bool {
	return recs.clients[clientId] != nil
}

func (recs *Records) Add(clientId string, rec *Record) (id string) {
	recs.RLock()
	client := recs.clients[clientId]
	recs.RUnlock()

	return client.Add(rec)
}

func (recs *Records) Remove(clientId, recordId string) *Record {
	recs.RLock()
	client := recs.clients[clientId]
	recs.RUnlock()

	if client == nil {
		return nil
	}

	return client.Remove(recordId)
}

func (recs *Records) RemoveClient(clientId string) {
	recs.Lock()
	defer recs.Unlock()

	close(recs.clients[clientId].Tunnel)
	recs.clients[clientId] = nil
}

func (recs *Records) Get(clientId string, bearer string) (*Client, error) {
	recs.RLock()
	defer recs.RUnlock()

	client := recs.clients[clientId]

	if client.Authorize(bearer) {
		return client, nil
	}

	return nil, errors.New("Authentication failed for client '" + clientId + "'")
}

func (recs *Record) ToString(maxPathLength int) string {
	path := []byte(recs.Request.Path)

	if len(path) > maxPathLength {
		path = append([]byte("..."), path[len(path)-maxPathLength:]...)
	}

	method := "[" + recs.Request.Method + "]"

	strPathLength := strconv.Itoa(maxPathLength + 3)
	str := fmt.Sprintf("%8s %"+strPathLength+"s", method, string(path))

	if recs.Response == nil {
		return str
	}

	return str + fmt.Sprintf(" %d - %dms", recs.Response.Code, recs.Elapsed)
}
