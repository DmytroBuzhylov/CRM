package usecase

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/feature/staff/interface_adapters/dto"
	"Test/internal/feature/staff/interface_adapters/mapper"
	"Test/internal/feature/staff/repository"
	"Test/internal/pkg/provider"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StaffUCInteractor struct {
	staffRepo           repository.StaffRepository
	workDayRepo         repository.WorkDayRepository
	breakRepo           repository.BreakRepository
	staffInvitationRepo repository.InvitationStaffRepository
	pool                *pgxpool.Pool
}

func NewStaffUCInteractor(
	staffRepo repository.StaffRepository,
	workDayRepo repository.WorkDayRepository,
	breakRepo repository.BreakRepository,
	staffInvitation repository.InvitationStaffRepository,
) *StaffUCInteractor {
	return &StaffUCInteractor{
		staffRepo:           staffRepo,
		workDayRepo:         workDayRepo,
		breakRepo:           breakRepo,
		staffInvitationRepo: staffInvitation,
	}
}

func (i *StaffUCInteractor) CreateInviteToStaff(ctx context.Context, staff *dto.InviteStaffReq) error {
	staffInvite := &entity.StaffInvitation{
		ID:             uuid.New(),
		OrganizationID: staff.OrganizationID,
		InvitedEmail:   staff.InvitedEmail,
		InviterUserID:  staff.InviterUserId,
		InvitationCode: uuid.New().String(),
		Position:       staff.Position,
		RoleID:         staff.RoleID,
		ExpiresAt:      time.Now().Add(60 * time.Minute),
		Status:         "pending",
		CreatedAt:      time.Now(),
	}
	err := i.staffInvitationRepo.Create(ctx, staffInvite)

	return err
}

func (i *StaffUCInteractor) AcceptInviteToStaff(ctx context.Context, userID uuid.UUID, inviteCode string) error {
	dataProvider := provider.NewPostgresProvider(i.pool)

	tx, err := dataProvider.BeginTx(ctx)
	if err != nil {
		return err
	}

	staffInvitationRepo := repository.NewInvitationStaffRepo(tx)
	staffRepo := repository.NewStaffRepo(dataProvider)

	staffInvite, err := staffInvitationRepo.GetByCode(ctx, inviteCode)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	staff := &entity.Staff{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: staffInvite.OrganizationID,
		RoleID:         staffInvite.RoleID,
		Position:       staffInvite.Position,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = staffRepo.Create(ctx, staff)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = staffInvitationRepo.Expire(ctx, staffInvite.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (i *StaffUCInteractor) DismissStaff(ctx context.Context, id uuid.UUID) error {
	return i.staffRepo.Delete(ctx, id)
}

func (i *StaffUCInteractor) AddTimetableStaff(ctx context.Context, req *dto.AddTimetableStaffReq) error {

	if req.StaffID == uuid.Nil || len(req.WorkDays) == 0 {
		return fmt.Errorf("validation error: staff ID and work days are required")
	}

	workDays, breaksByWorkDayID := mapper.MapAddTimetableReqToEntities(req)

	if err := i.staffRepo.SaveTimetable(ctx, req.StaffID, workDays, breaksByWorkDayID); err != nil {
		return fmt.Errorf("failed to save staff timetable: %w", err)
	}

	return nil
}

func (i *StaffUCInteractor) DeleteTimetableStaff(ctx context.Context, staffID uuid.UUID, days []time.Weekday) error {
	return i.staffRepo.DeleteTimetable(ctx, staffID, days)
}
