package usecase

import (
	"Test/internal/feature/staff/interface_adapters/dto"
	"context"
	"github.com/google/uuid"
	"time"
)

type StaffUseCase interface {
	CreateInviteToStaff(ctx context.Context, staff *dto.InviteStaffReq) error
	AcceptInviteToStaff(ctx context.Context, userID uuid.UUID, inviteCode string) error
	DismissStaff(ctx context.Context, id uuid.UUID) error
	AddTimetableStaff(ctx context.Context, req *dto.AddTimetableStaffReq) error
	DeleteTimetableStaff(ctx context.Context, staffID uuid.UUID, days []time.Weekday) error
}
