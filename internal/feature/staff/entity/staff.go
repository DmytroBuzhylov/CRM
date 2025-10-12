package entity

import (
	"github.com/google/uuid"
	"time"
)

type Staff struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	RoleID         uuid.UUID
	Position       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
