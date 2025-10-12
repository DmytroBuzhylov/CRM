package entity

import (
	"github.com/google/uuid"
	"time"
)

type WorkDay struct {
	ID         uuid.UUID
	StaffID    uuid.UUID
	DayOfWeak  time.Weekday
	ShiftStart time.Time
	ShiftEnd   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
