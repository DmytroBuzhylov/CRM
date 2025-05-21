package entity

import "time"

type Note struct {
	ID        uint64
	AuthorID  uint64
	ClientID  uint64
	Text      string
	CreatedAt time.Time
}
