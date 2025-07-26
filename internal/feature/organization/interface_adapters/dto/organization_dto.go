package dto

type CreateOrganizationRequest struct {
	OwnerUserID uint64 `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateOrganizationResponse struct {
	Status string `json:"status"`
}

type CreateInviteRequest struct {
	OrganizationID *uint64 `json:"-"` //from jwt claims
	InvitedEmail   string  `json:"invited_email" binding:"required,email"`
}

type CreateInviteResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
}

type AcceptInvitationRequest struct {
	UserID uint64 `json:"-"`
	Code   string `json:"-"`
}

type AcceptInvitationResponse struct {
	Status string `json:"status"`
}
