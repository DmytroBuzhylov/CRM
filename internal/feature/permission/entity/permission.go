package entity

import (
	"github.com/google/uuid"
	"time"
)

type Permission struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
