package usecase

import (
	"Test/internal/feature/notification/entity"
	notification_const "Test/internal/feature/notification/entity/const"
	"Test/internal/feature/notification/interface_adapters/dto"
	"Test/internal/feature/notification/interface_adapters/email_sender"
	"Test/internal/feature/notification/interface_adapters/mapper"
	"Test/internal/feature/notification/repository"
	repository2 "Test/internal/feature/user/repository"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type NotificationUseCaseInteractor struct {
	repo     repository.NotificationRepository
	userRepo repository2.UserRepository
	sender   email_sender.Sender
}

func NewNotificationUseCaseInteractor(repo repository.NotificationRepository, userRepo repository2.UserRepository, sender email_sender.Sender) *NotificationUseCaseInteractor {
	return &NotificationUseCaseInteractor{repo: repo, sender: sender, userRepo: userRepo}
}

func (uc *NotificationUseCaseInteractor) CreateAndDeliver(ctx context.Context, req *entity.Notification) (*entity.Notification, error) {
	req.ID = uuid.New()
	req.CreatedAt = time.Now().Unix()

	if req.SendToEmail {
		err := func() error {
			user, err := uc.userRepo.FindById(ctx, req.UserID)
			if err != nil {
				return errors.New("find user email error")
			}

			templateName, subject, err := mapper.PrepareEmailContent(req.Type)
			if err != nil {
				return err
			}
			data := struct {
				ResourceType      notification_const.ResourceType
				RelatedResourceID string
				Text              string
			}{
				RelatedResourceID: req.RelatedResourceID.String(),
				ResourceType:      *req.ResourceType,
				Text:              req.Text,
			}

			err = uc.sender.Send(user.Email, subject, templateName, data)
			if err != nil {
				fmt.Println(err)
				return errors.New("error send email")
			}
			return nil
		}()
		if err != nil {
			log.Warn().
				Str("usecase", "NotificationUseCaseInteractor").
				Str("method", "CreateAndDeliver").
				Str("error", "failed to send email").
				Err(err).
				Send()
		}
	}
	err := uc.repo.Create(ctx, req)
	if err != nil {
		log.Warn().
			Str("usecase", "NotificationUseCaseInteractor").
			Str("method", "CreateAndDeliver").
			Err(err).
			Send()

		return nil, errors.New("notification save error")
	}

	return req, nil
}

func (uc *NotificationUseCaseInteractor) ViewNotification(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (uc *NotificationUseCaseInteractor) GetUserNotification(ctx context.Context, notificationID uuid.UUID) (*entity.Notification, error) {
	return nil, nil
}

func (uc *NotificationUseCaseInteractor) GetUserNotifications(ctx context.Context, userID uuid.UUID) ([]*dto.GetNotificationResponse, error) {

	notifications, err := uc.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var resp []*dto.GetNotificationResponse
	for _, n := range notifications {
		resType := string(*n.ResourceType)
		resp = append(resp, &dto.GetNotificationResponse{
			ID:                n.ID,
			OrganizationID:    n.OrganizationID,
			UserID:            n.UserID,
			Type:              string(n.Type),
			Text:              n.Text,
			Viewed:            n.Viewed,
			RelatedResourceID: n.RelatedResourceID,
			ResourceType:      resType,
			CreatedAt:         n.CreatedAt,
		})
	}

	return resp, nil
}

func (uc *NotificationUseCaseInteractor) GetOrganizationNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	return nil, nil
}
