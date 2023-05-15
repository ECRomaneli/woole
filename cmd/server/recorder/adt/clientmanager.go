package adt

import (
	"bytes"
	"encoding/hex"
	"errors"
	"sync"
	"woole/shared/util/rand"
)

type ClientManager struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{clients: make(map[string]*Client)}
}

func (cm *ClientManager) Register(clientId string, bearer []byte) *Client {
	clientId = cm.generateClientId(clientId)

	client := NewClient(clientId, bearer)
	cm.put(clientId, client)

	return client
}

func (cm *ClientManager) Deregister(clientId string) {
	close(cm.clients[clientId].NewRecordChannel)
	cm.put(clientId, nil)
}

func (cm *ClientManager) DeregisterIfIdle(clientId string, callback func()) {
	client := cm.clients[clientId]

	go func() {
		<-client.IdleTimeout.C
		cm.Deregister(client.Id)
		callback()
	}()
}

func (cm *ClientManager) RecoverSession(clientId string, bearer []byte) (*Client, error) {
	client := cm.Get(clientId)

	if client == nil {
		return nil, nil
	}

	if !bytes.Equal(client.Bearer, bearer) {
		return nil, errors.New("failed to authenticate the client")
	}

	return client, nil
}

func (cm *ClientManager) Get(clientId string) *Client {
	return cm.clients[clientId]
}

func (cm *ClientManager) Exists(clientId string) bool {
	return cm.clients[clientId] != nil
}

func (cm *ClientManager) generateClientId(clientId string) string {
	hasClientId := clientId != ""

	if !hasClientId {
		return hex.EncodeToString(rand.RandMD5(""))[:8]
	}

	if cm.Exists(clientId) {
		return clientId + "-" + hex.EncodeToString(rand.RandMD5(clientId))[:5]
	}

	return clientId
}

func (cm *ClientManager) put(clientId string, client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[clientId] = client
}
