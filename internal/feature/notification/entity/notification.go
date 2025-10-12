package entity

import (
	notification_const "Test/internal/feature/notification/entity/const"
	"github.com/google/uuid"
)

type Notification struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID

	Type        notification_const.NotificationType
	Text        string
	SendToEmail bool
	SendToApp   bool
	Viewed      bool

	RelatedResourceID *uuid.UUID
	ResourceType      *notification_const.ResourceType

	CreatedAt int64
}
