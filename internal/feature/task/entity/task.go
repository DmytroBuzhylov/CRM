package entity

import (
	"github.com/google/uuid"
	"github.com/maniartech/signals"
	"time"
)

var (
	TaskCreatedSignal         = signals.New[*Task]()
	TaskUpdatedSignal         = signals.New[*Task]()
	TaskUpdatedDeadlineSignal = signals.New[*Task]()
	TaskUpdatedStatusSignal   = signals.New[*Task]()
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
