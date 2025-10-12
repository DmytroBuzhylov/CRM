package event_listener

import (
	"context"
)

type Listener interface {
	Listen(ctx context.Context)
	Close() error
}
