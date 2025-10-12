package repository

import (
	"Test/internal/feature/notification/entity"
	"context"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, notify *entity.Notification) error
	UpdateViewStatus(ctx context.Context, id uuid.UUID) error

	Get(ctx context.Context, id uuid.UUID) (*entity.Notification, error)
	GetMany(ctx context.Context, organizationID uuid.UUID) ([]*entity.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error)
}
