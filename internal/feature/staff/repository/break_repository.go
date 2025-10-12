package repository

import (
	"Test/internal/feature/staff/entity"
	"context"
	"github.com/google/uuid"
)

type BreakRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.Break, error)
	GetMany(ctx context.Context, workDayId uuid.UUID) ([]*entity.Break, error)
	Create(ctx context.Context, br *entity.Break) error
	Update(ctx context.Context, br *entity.Break) error
	Delete(ctx context.Context, id uuid.UUID) error
}
