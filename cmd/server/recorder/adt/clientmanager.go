package adt

import (
	"bytes"
	"encoding/hex"
	"errors"
	"sync"
	pb "woole/shared/payload"
	"woole/shared/util/hash"
)

type ClientManager struct {
	mu      sync.Mutex
	clients map[string]*Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{clients: make(map[string]*Client)}
}

func (cm *ClientManager) Register(clientId string) *Client {
	clientId = cm.generateClientId(clientId)

	client := NewClient(clientId)
	cm.put(clientId, client)

	return client
}

func (cm *ClientManager) Deregister(clientId string) {
	close(cm.clients[clientId].recordChannel)
	cm.put(clientId, nil)
}

func (cm *ClientManager) RecoverSession(session *pb.Session) (*Client, error) {
	client := cm.Get(session.ClientId)

	if client == nil {
		return nil, errors.New("The client '" + session.ClientId + "' is not in use")
	}

	if !bytes.Equal(client.Bearer, session.Bearer) {
		return nil, errors.New("Failed to authenticate to client '" + session.ClientId + "'")
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

	for clientId == "" || cm.Exists(clientId) {
		if hasClientId {
			clientId = clientId + "-" + hex.EncodeToString(hash.RandMD5(clientId))[:5]
		} else {
			clientId = hex.EncodeToString(hash.RandMD5(""))[:8]
		}
	}

	return clientId
}

func (cm *ClientManager) put(clientId string, client *Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[clientId] = client
}
