package repository

import (
	"Test/internal/feature/task/entity"
	"context"
	"github.com/google/uuid"
	"time"
)

type TaskCreator interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
}

type TaskGetter interface {
	Get(ctx context.Context, id uuid.UUID) (entity.Task, error)
	GetTasks(ctx context.Context, filter TaskFilter, pagination Pagination) ([]entity.Task, error)
}

type TaskUpdater interface {
	Update(ctx context.Context, task entity.Task) (entity.Task, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateDeadline(ctx context.Context, id uuid.UUID, deadline time.Time) error
}

type TaskDeleter interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type TaskFilter struct {
	Name       string
	Priority   uint
	Status     string
	Deadline   *time.Time
	AssigneeID uuid.UUID
	ClientID   uuid.UUID
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type Pagination struct {
	Page  uint
	Limit uint
}
