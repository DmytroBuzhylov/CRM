package repository

import (
	"Test/internal/feature/user/entity"
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	FindById(ctx context.Context, id uint64) (entity.User, error)
	SaveRefreshToken(ctx context.Context, userID uint64, tokenID string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, tokenID string) error
	FindRefreshToken(ctx context.Context, tokenID string) (uint64, error)
}
