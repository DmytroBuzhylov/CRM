package ws_hub

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type gorillaWebsocket struct {
	upgrader websocket.Upgrader
	hub      *Hub
}

func NewGorillaWebsocket(upgrader websocket.Upgrader) (*gorillaWebsocket, Broadcaster) {

	clients := NewInMemoryClients()
	regCh := make(chan *Client, 1024)
	unregCh := make(chan *Client, 1024)
	mesCh := make(chan *HubMessage, 1024)

	hub := newHub(clients, mesCh, regCh, unregCh)

	go hub.Run()

	gWs := &gorillaWebsocket{
		upgrader: upgrader,
		hub:      hub,
	}

	broadcaster := NewWsBroadcaster(mesCh)

	return gWs, broadcaster
}

func (ws *gorillaWebsocket) Handle(c *gin.Context) {
	stringUserID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user id"})
		return
	}
	userID, ok := stringUserID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect user id format"})
		return
	}

	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := newClient(ws.hub, conn, userID)
	ws.hub.register <- client

	go client.writePump()

	client.readPump()
}
