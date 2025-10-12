package dto

import "github.com/google/uuid"

//	type SendNotificationRequest struct {
//		OrganizationID uuid.UUID
//		UserID         uuid.UUID
//
//		Type        entity.NotificationType
//		Text        string
//		SendToEmail bool
//		SendToApp   bool
//		Viewed      bool
//
//		RelatedResourceID *uuid.UUID
//		ResourceType      *entity.ResourceType
//	}
type GetNotificationResponse struct {
	ID             uuid.UUID `json:"id,omitempty"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	UserID         uuid.UUID `json:"user_id,omitempty"`

	Type   string `json:"type,omitempty"`
	Text   string `json:"text,omitempty"`
	Viewed bool   `json:"viewed"`

	RelatedResourceID *uuid.UUID `json:"related_resource_id,omitempty"`
	ResourceType      string     `json:"resource_type,omitempty"`

	CreatedAt int64 `json:"created_at,omitempty"`
}
