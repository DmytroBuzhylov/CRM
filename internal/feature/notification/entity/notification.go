package entity

import (
	notification_const2 "Test/internal/pkg/const"
	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID

	Type        notification_const2.NotificationType
	Text        string
	SendToEmail bool
	SendToApp   bool
	Viewed      bool

	RelatedResourceID *uuid.UUID
	ResourceType      *notification_const2.ResourceType

	CreatedAt int64
}
