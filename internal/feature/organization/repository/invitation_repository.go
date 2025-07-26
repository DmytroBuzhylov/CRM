package repository

import (
	"Test/internal/feature/organization/entity"
	"context"
)

type InvitationRepository interface {
	Save(ctx context.Context, inv *entity.Invitation) error
	GetByCode(ctx context.Context, code string) (*entity.Invitation, error)
	UpdateStatus(ctx context.Context, id uint64, status string) error
	AcceptInvite(ctx context.Context, userID uint64, code string) error
}
