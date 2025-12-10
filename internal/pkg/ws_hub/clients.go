package ws_hub

import (
	"github.com/google/uuid"
	"sync"
)

type Clients interface {
	SetConn(client *Client) error
	GetConn(userID uuid.UUID) ([]*Client, error)
	Delete(client *Client) error
}

type InMemoryClients struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[*Client]bool
}

func NewInMemoryClients() *InMemoryClients {
	return &InMemoryClients{
		clients: make(map[uuid.UUID]map[*Client]bool),
	}
}

func (c *InMemoryClients) SetConn(client *Client) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.clients[client.UserID]; !ok {
		c.clients[client.UserID] = make(map[*Client]bool)
	}
	c.clients[client.UserID][client] = true
	return nil
}

func (c *InMemoryClients) GetConn(userID uuid.UUID) ([]*Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	connMap, ok := c.clients[userID]
	if !ok {
		return []*Client{}, nil
	}

	if len(connMap) == 0 {
		return []*Client{}, nil
	}

	conns := make([]*Client, 0, len(connMap))
	for k := range connMap {
		conns = append(conns, k)
	}
	return conns, nil
}

func (c *InMemoryClients) Delete(client *Client) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if connMap, ok := c.clients[client.UserID]; ok {
		if _, ok := connMap[client]; ok {
			delete(c.clients[client.UserID], client)
			if len(c.clients[client.UserID]) == 0 {
				delete(c.clients, client.UserID)
			}
		}
	}
	return nil
}
