package events

import "time"

type EventBus interface {
	SendEvent(ev *Event)
}

type ChanEventBus struct {
	Events chan<- *Event
	Errors chan<- *EventError
}

func NewEventBus(bufferSize int) (*EventBus, <-chan *Event, <-chan *EventError) {
	evCh := make(chan *Event, bufferSize)
	errCh := make(chan *EventError, bufferSize)
	return &EventBus{Events: evCh, Errors: errCh}, evCh, errCh
}

type Event struct {
	Data    []byte
	Source  string
	Time    time.Time
	Headers map[string]string
}

type EventError struct {
	err        error
	ResourceID string
}

func (e *EventError) Error() string {
	return e.err.Error()
}

func NewEventError(err error, resourceId string) *EventError {
	return &EventError{
		err:        err,
		ResourceID: resourceId,
	}
}
