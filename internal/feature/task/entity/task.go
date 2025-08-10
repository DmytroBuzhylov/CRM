package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID             uuid.UUID
	Name           string
	Description    string
	Priority       uint
	Status         string
	Deadline       *time.Time
	AssigneeID     uuid.UUID
	ClientID       uuid.UUID
	OrganizationID uuid.UUID
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
