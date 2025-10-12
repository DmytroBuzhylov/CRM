package repository

import (
	"Test/internal/feature/staff/entity"
	"context"
	"github.com/google/uuid"
)

type WorkDayRepository interface {
	Create(ctx context.Context, wD *entity.WorkDay) error
	Get(ctx context.Context, id uuid.UUID) (*entity.WorkDay, error)
	GetMany(ctx context.Context, staffId uuid.UUID) ([]*entity.WorkDay, error)
	Update(ctx context.Context, wD *entity.WorkDay) error
	Delete(ctx context.Context, id uuid.UUID) error
}
