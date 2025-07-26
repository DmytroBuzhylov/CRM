package entity

import "time"

type Invitation struct {
	ID             uint64    `json:"id"`
	OrganizationID uint64    `json:"organization_id"`
	InvitedEmail   string    `json:"invited_email"`
	InvitationCode string    `json:"invitation_code"`
	ExpiresAt      time.Time `json:"expires_at"`
	Status         string    `json:"status"` // "pending", "accepted", "declined", "expired"
	CreatedAt      time.Time `json:"created_at"`
}
