package repository

import (
	"Test/internal/feature/table/entity"
	"context"
	"github.com/google/uuid"
)

type TableRepository interface {
	Create(ctx context.Context, table *entity.Table) error
	UpdateSeats(ctx context.Context, seats uint64, id uuid.UUID, organizationID uuid.UUID) error
	UpdateNumber(ctx context.Context, number uint64, id uuid.UUID, organizationID uuid.UUID) error
	UpdateStatus(ctx context.Context, status string, id uuid.UUID, organizationID uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Table, error)
	GetAll(ctx context.Context, organizationID uuid.UUID) ([]*entity.Table, error)
	Delete(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) error
}
