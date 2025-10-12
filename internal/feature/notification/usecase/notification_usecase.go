package usecase

import (
	"Test/internal/feature/notification/entity"
	"Test/internal/feature/notification/interface_adapters/dto"
	"context"
	"github.com/google/uuid"
)

type NotificationUseCase interface {
	CreateAndDeliver(ctx context.Context, req *entity.Notification) (*entity.Notification, error)

	ViewNotification(ctx context.Context, id uuid.UUID) error

	GetUserNotification(ctx context.Context, notificationID uuid.UUID) (*entity.Notification, error)
	GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]*dto.GetNotificationResponse, error)

	GetOrganizationNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
}
