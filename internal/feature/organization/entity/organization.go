package entity

import "time"

type Organization struct {
	ID          uint64
	Name        string
	Description string
	OwnerUserID uint64
	CreatedAt   time.Time
}
