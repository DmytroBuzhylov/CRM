package repository

import (
	"Test/internal/feature/staff/entity"
	"context"
	"github.com/google/uuid"
)

type InvitationStaffRepository interface {
	Get(ctx context.Context, id uuid.UUID) (*entity.StaffInvitation, error)
	GetByCode(ctx context.Context, code string) (*entity.StaffInvitation, error)
	Create(ctx context.Context, invite *entity.StaffInvitation) error
	Expire(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
