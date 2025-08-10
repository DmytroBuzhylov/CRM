package repository

import (
	"Test/internal/feature/user/entity"
	"context"
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindByPhone(ctx context.Context, phone string) (entity.User, error)
	FindById(ctx context.Context, id uuid.UUID) (entity.User, error)
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiresAt time.Time) error
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	FindRefreshToken(ctx context.Context, tokenID string) (RefreshToken, error)
	GetOrganizationID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}

type RefreshToken struct {
	ID        string
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt time.Time
}
