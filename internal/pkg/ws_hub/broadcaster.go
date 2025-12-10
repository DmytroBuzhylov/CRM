package ws_hub

import "github.com/google/uuid"

type Broadcaster interface {
	SendMessage(userID uuid.UUID, topic string, payload []byte) error
}

type wsBroadcaster struct {
	mesCh chan *HubMessage
}

func NewWsBroadcaster(mesCh chan *HubMessage) *wsBroadcaster {
	return &wsBroadcaster{mesCh: mesCh}
}

func (b *wsBroadcaster) SendMessage(userID uuid.UUID, topic string, payload []byte) error {
	b.mesCh <- &HubMessage{
		UserID:  userID,
		Topic:   topic,
		Payload: payload,
	}
	return nil
}
