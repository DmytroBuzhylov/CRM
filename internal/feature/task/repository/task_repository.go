package repository

import (
	"Test/internal/feature/task/entity"
	"context"
	"time"
)

type TaskCreator interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
}

type TaskGetter interface {
	Get(ctx context.Context, id uint64) (entity.Task, error)
	GetTasks(ctx context.Context, filter TaskFilter, pagination Pagination) ([]entity.Task, error)
}

type TaskUpdater interface {
	Update(ctx context.Context, task entity.Task) (entity.Task, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
	UpdateDeadline(ctx context.Context, id uint64, deadline time.Time) error
}

type TaskDeleter interface {
	Delete(ctx context.Context, id uint64) error
}

type TaskRepository interface {
	TaskCreator
	TaskGetter
	TaskUpdater
	TaskDeleter
}

type TaskFilter struct {
	Name       string
	Priority   uint
	Status     string
	Deadline   *time.Time
	AssigneeID uint64
	ClientID   uint64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

type Pagination struct {
	Page  uint
	Limit uint
}
