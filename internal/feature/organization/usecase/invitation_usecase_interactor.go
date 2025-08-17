package usecase

import (
	"Test/internal/feature/organization/entity"
	"Test/internal/feature/organization/interface_adapters/dto"
	"Test/internal/feature/organization/repository"
	"Test/pkg/utils"
	"context"
	"time"
)

type invitationUseCaseInteractor struct {
	repo repository.InvitationRepository
}

func NewInvitationUseCaseInteractor(repo repository.InvitationRepository) *invitationUseCaseInteractor {
	return &invitationUseCaseInteractor{repo: repo}
}

func (uc *invitationUseCaseInteractor) GenerateInvitation(ctx context.Context, req dto.CreateInviteRequest) (dto.CreateInviteResponse, error) {
	code := utils.GenerateInvitationCode()
	newInvitation := entity.Invitation{
		OrganizationID: req.OrganizationID,
		InvitedEmail:   req.InvitedEmail,
		InvitationCode: code,
		ExpiresAt:      time.Now().Add(time.Hour * 24 * 7),
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	err := uc.repo.Save(ctx, newInvitation)
	if err != nil {
		return dto.CreateInviteResponse{}, err
	}

	return dto.CreateInviteResponse{Status: "ok", Code: code}, nil
}

func (uc *invitationUseCaseInteractor) AcceptInvitation(ctx context.Context, req dto.AcceptInvitationRequest) (dto.AcceptInvitationResponse, error) {
	err := uc.repo.AcceptInvite(ctx, req.UserID, req.Code)
	if err != nil {
		return dto.AcceptInvitationResponse{}, err
	}

	return dto.AcceptInvitationResponse{Status: "ok"}, nil
}
