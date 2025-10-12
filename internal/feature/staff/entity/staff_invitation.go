package entity

import (
	"github.com/google/uuid"
	"time"
)

type StaffInvitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	InvitedEmail   string
	InviterUserId  uuid.UUID
	InvitationCode string
	Position       string
	RoleID         uuid.UUID
	ExpiresAt      time.Time
	Status         string
	CreatedAt      time.Time
}
