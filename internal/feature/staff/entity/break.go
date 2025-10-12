package entity

import (
	"github.com/google/uuid"
	"time"
)

type Break struct {
	ID         uuid.UUID
	WorkDayID  uuid.UUID
	BreakBegin time.Time
	BreakEnd   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
