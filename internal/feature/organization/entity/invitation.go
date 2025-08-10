package entity

import (
	"github.com/google/uuid"
	"time"
)

type Invitation struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	InvitedEmail   string    `json:"invited_email"`
	InvitationCode string    `json:"invitation_code"`
	ExpiresAt      time.Time `json:"expires_at"`
	Status         string    `json:"status"` // "pending", "accepted", "declined", "expired"
	CreatedAt      time.Time `json:"created_at"`
}
