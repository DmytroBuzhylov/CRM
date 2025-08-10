package dto

import "github.com/google/uuid"

type CreateOrganizationRequest struct {
	OwnerUserID uuid.UUID `json:"-"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type CreateOrganizationResponse struct {
	Status string `json:"status"`
}

type CreateInviteRequest struct {
	OrganizationID uuid.UUID `json:"-"` //from jwt claims
	InvitedEmail   string    `json:"invited_email" binding:"required,email"`
}

type CreateInviteResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
}

type AcceptInvitationRequest struct {
	UserID uuid.UUID `json:"-"`
	Code   string    `json:"-"`
}

type AcceptInvitationResponse struct {
	Status string `json:"status"`
}
