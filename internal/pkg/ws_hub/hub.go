package ws_hub

import (
	"github.com/google/uuid"
	"log"
)

type HubMessage struct {
	UserID  uuid.UUID
	Topic   string
	Payload []byte
}

type Hub struct {
	clients Clients

	register chan *Client

	unregister chan *Client

	broadcast chan *HubMessage
}

func newHub(clients Clients, broadcast chan *HubMessage, register chan *Client, unregister chan *Client) *Hub {
	return &Hub{
		clients:    clients,
		register:   register,
		unregister: unregister,
		broadcast:  broadcast,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if err := h.clients.SetConn(client); err != nil {
				log.Printf("Failed to set connection for client %s: %v", client.UserID, err)
			} else {
				log.Printf("Client registered: %s", client.UserID)
			}

		case client := <-h.unregister:
			if err := h.clients.Delete(client); err != nil {
				log.Printf("Failed to delete connection for client %s: %v", client.UserID, err)
			} else {
				log.Printf("Client unregistered: %s", client.UserID)
			}
			close(client.send)

		case message := <-h.broadcast:
			clients, err := h.clients.GetConn(message.UserID)
			if err != nil {
				continue
			}

			for _, client := range clients {
				select {
				case client.send <- message.Payload:
				default:
					log.Printf("Client send buffer full. Closing client: %s", client.UserID)
					go func(c *Client) {
						h.unregister <- c
					}(client)
				}
			}
		}
	}
}
