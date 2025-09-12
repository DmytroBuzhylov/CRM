package entity

import (
	"github.com/google/uuid"
	"time"
)

type Table struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Number         uint64
	Status         string
	Seats          uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
