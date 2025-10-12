package repository

import (
	"Test/internal/feature/staff/entity"
	"context"
	"github.com/google/uuid"
	"time"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *entity.Staff) error
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Staff, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.Staff, error)
	GetMany(ctx context.Context, organizationId uuid.UUID) ([]*entity.Staff, error)
	GetRoleIDByUserID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)

	SaveTimetable(ctx context.Context, staffID uuid.UUID, workDays []entity.WorkDay, breaks map[uuid.UUID][]entity.Break) error
	DeleteTimetable(ctx context.Context, staffID uuid.UUID, days []time.Weekday) error
}
