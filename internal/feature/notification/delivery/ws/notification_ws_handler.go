package ws

import (
	"Test/internal/pkg/events"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	UserID uuid.UUID
	Data   []byte
}

type Client struct {
	userID uuid.UUID
	conn   *websocket.Conn
}

type Hub struct {
	clients    map[uuid.UUID]map[*websocket.Conn]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	evChan     <-chan *events.Event
	errChan    <-chan *events.EventError
	mu         sync.Mutex
}

type NotificationWSHandler struct {
	upgrader websocket.Upgrader
	hub      Hub
}

func NewNotificationWSHandler(
	upgrader websocket.Upgrader,
	evChan <-chan *events.Event,
	errChan <-chan *events.EventError,
) *NotificationWSHandler {
	nh := &NotificationWSHandler{
		upgrader: upgrader,
		hub: Hub{
			clients:    make(map[uuid.UUID]map[*websocket.Conn]bool),
			register:   make(chan *Client, 128),
			unregister: make(chan *Client, 128),
			broadcast:  make(chan *Message, 1024),
			evChan:     evChan,
			errChan:    errChan,
		},
	}

	go nh.hub.Run()

	return nh
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*websocket.Conn]bool)
			}
			h.clients[client.userID][client.conn] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			delete(h.clients[client.userID], client.conn)
			h.mu.Unlock()

		case msg := <-h.broadcast:
			for conn, ok := range h.clients[msg.UserID] {
				if !ok {
					continue
				}
				err := conn.WriteMessage(websocket.BinaryMessage, msg.Data)
				if err != nil {
					continue
				}
			}
		case msg := <-h.evChan:
			userId, _ := uuid.Parse(msg.Headers["user_id"])
			h.broadcast <- &Message{
				UserID: userId,
				Data:   msg.Data,
			}
		case evErr := <-h.errChan:
			log.Warn().Str("resource_id", evErr.ResourceID).
				Str("error", evErr.Error()).
				Send()
		}
	}
}

func (h *NotificationWSHandler) NotificationWSHandler(c *gin.Context) {

	stringUserID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "incorrect user id",
		})
		return
	}
	userID := stringUserID.(uuid.UUID)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "websocket handshake error",
		})
		return
	}

	client := &Client{
		userID: userID,
		conn:   conn,
	}

	h.hub.register <- client
	defer func() {
		h.hub.unregister <- client
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	ticker := time.NewTicker(55 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
				return
			}
		}
	}

}
