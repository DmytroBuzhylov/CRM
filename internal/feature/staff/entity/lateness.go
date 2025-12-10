package entity

import (
	"github.com/google/uuid"
	"time"
)

type Lateness struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	StaffID        uuid.UUID
	TimeOfDelay    time.Time
	CreatedAt      time.Time
}
