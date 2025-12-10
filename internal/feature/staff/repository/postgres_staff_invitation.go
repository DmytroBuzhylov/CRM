package repository

import (
	"Test/internal/feature/staff/entity"
	"Test/internal/pkg/provider"
	"context"
	_ "embed"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

//go:embed sql/staff_invitation/create.sql
var createStaffInvitationQuery string

//go:embed sql/staff_invitation/get_by_id.sql
var getStaffInvitationByIDQuery string

//go:embed sql/staff_invitation/get_by_code.sql
var getStaffInvitationByCodeQuery string

//go:embed sql/staff_invitation/update_expire.sql
var updateStaffInvitationExpireQuery string

//go:embed sql/staff_invitation/update_status.sql
var updateStaffInvitationStatusQuery string

type InvitationStaffRepo struct {
	connector provider.SQLConnector
	qb        sq.StatementBuilderType
}

func NewInvitationStaffRepo(connector provider.SQLConnector) *InvitationStaffRepo {
	return &InvitationStaffRepo{
		connector: connector,
		qb:        sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *InvitationStaffRepo) Create(ctx context.Context, invite *entity.StaffInvitation) error {
	_, err := r.connector.Exec(ctx, createStaffInvitationQuery,
		invite.ID,
		invite.OrganizationID,
		invite.InvitedEmail,
		invite.InviterUserID,
		invite.InvitationCode,
		invite.Position,
		invite.RoleID,
		invite.ExpiresAt,
		invite.Status,
		invite.CreatedAt,
	)
	return err
}

func (r *InvitationStaffRepo) Expire(ctx context.Context, id uuid.UUID) error {

	_, err := r.connector.Exec(ctx, updateStaffInvitationExpireQuery, id)
	return err
}

func (r *InvitationStaffRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {

	_, err := r.connector.Exec(ctx, updateStaffInvitationStatusQuery, status, id)
	return err
}

func (r *InvitationStaffRepo) Get(ctx context.Context, id uuid.UUID) (*entity.StaffInvitation, error) {
	var invite = new(entity.StaffInvitation)

	err := r.connector.QueryRow(ctx, getStaffInvitationByIDQuery).Scan(
		&invite.OrganizationID,
		&invite.InvitedEmail,
		&invite.InviterUserID,
		&invite.InvitationCode,
		&invite.Position,
		&invite.RoleID,
		&invite.ExpiresAt,
		&invite.Status,
		&invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	invite.ID = id

	return invite, nil
}

func (r *InvitationStaffRepo) GetByCode(ctx context.Context, code string) (*entity.StaffInvitation, error) {
	var invite = new(entity.StaffInvitation)

	err := r.connector.QueryRow(ctx, getStaffInvitationByCodeQuery).Scan(
		&invite.ID,
		&invite.OrganizationID,
		&invite.InvitedEmail,
		&invite.InviterUserID,
		&invite.InvitationCode,
		&invite.Position,
		&invite.RoleID,
		&invite.ExpiresAt,
		&invite.Status,
		&invite.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return invite, nil
}
