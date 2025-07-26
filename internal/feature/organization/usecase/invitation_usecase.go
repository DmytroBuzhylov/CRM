package usecase

import (
	"Test/internal/feature/organization/interface_adapters/dto"
	"context"
)

type InvitationUseCase interface {
	GenerateInvitation(ctx context.Context, req dto.CreateInviteRequest) (dto.CreateInviteResponse, error)
	AcceptInvitation(ctx context.Context, req dto.AcceptInvitationRequest) (dto.AcceptInvitationResponse, error)
}
