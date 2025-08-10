package entity

import (
	"github.com/google/uuid"
	"time"
)

type MenuItem struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    string
	Price          uint64
	Category       string
	IsAvailable    bool
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
