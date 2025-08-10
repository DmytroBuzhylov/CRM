package repository

import (
	"Test/internal/feature/organization/entity"
	"context"
	"github.com/google/uuid"
)

type InvitationRepository interface {
	Save(ctx context.Context, inv entity.Invitation) error
	GetByCode(ctx context.Context, code string) (entity.Invitation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	AcceptInvite(ctx context.Context, userID uuid.UUID, code string) error
}
