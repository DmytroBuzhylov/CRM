package entity

import (
	"github.com/google/uuid"
	"time"
)

type Role struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
