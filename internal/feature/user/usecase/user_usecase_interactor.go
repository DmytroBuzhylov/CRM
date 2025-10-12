package usecase

import (
	"Test/config"
	"Test/internal/feature/user/entity"
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"Test/internal/feature/user/repository/postgres"
	"Test/pkg/jwt"
	"Test/pkg/store"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

type UserUseCaseInteractor struct {
	repo      repository.UserRepository
	store     store.Store
	jwtConfig config.JWTConfig
}

func NewUserUseCaseInteractor(repo repository.UserRepository, jwtConfig config.JWTConfig, store store.Store) *UserUseCaseInteractor {
	return &UserUseCaseInteractor{repo: repo, jwtConfig: jwtConfig, store: store}
}

func (i *UserUseCaseInteractor) Create(ctx context.Context, req dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	hashedPassword, err := postgres.HashPassword(req.Password)
	if err != nil {
		return dto.CreateUserResponse{}, errors.New("error hashing password")
	}
	user := entity.NewUser(
		req.FullName,
		req.Username,
		hashedPassword,
		req.Email,
		req.Phone,
		"platform_user",
	)

	err = i.repo.Create(ctx, user)
	if err != nil {
		return dto.CreateUserResponse{}, err
	}

	newAccessToken, newRefreshToken, err := jwt.GenerateTokens(
		user.ID,
		user.Role,
		uuid.Nil,
		i.jwtConfig.JWTAccessSecret,
		i.jwtConfig.JWTRefreshSecret,
		i.jwtConfig.JWTAccessLifetime,
		i.jwtConfig.JWTRefreshLifetime,
	)
	if err != nil {
		return dto.CreateUserResponse{}, errors.New("error generate jwt tokens")
	}

	err = i.repo.SaveRefreshToken(ctx, user.ID, newRefreshToken, time.Now().Add(i.jwtConfig.JWTRefreshLifetime))
	if err != nil {
		_ = i.repo.DeleteUser(ctx, user.ID)
		return dto.CreateUserResponse{}, errors.New("error save jwt tokens")
	}

	return dto.CreateUserResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (i *UserUseCaseInteractor) GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *UserUseCaseInteractor) GetByUsername(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *UserUseCaseInteractor) GetByEmail(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *UserUseCaseInteractor) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var (
		user entity.User
		err  error
	)
	switch req.Method {
	case "email":
		user, err = i.repo.FindByEmail(ctx, req.Email)

	case "phone":
		user, err = i.repo.FindByPhone(ctx, req.Phone)

	case "username":
		user, err = i.repo.FindByUsername(ctx, req.Username)

	default:
		return dto.LoginResponse{Status: "incorrect method"}, fmt.Errorf("incorrect method")
	}

	if err != nil {
		return dto.LoginResponse{Status: "db error"}, errors.New("could not find username or phone or email")
	}
	if !postgres.CheckPasswordHash(req.Password, user.HashedPassword) {
		return dto.LoginResponse{Status: "incorrect password"}, errors.New("invalid password")
	}

	organizationID, err := i.repo.GetOrganizationID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			organizationID = uuid.Nil
		} else {
			return dto.LoginResponse{Status: "db error"}, errors.New("could not be found organization_id")
		}
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, organizationID, i.jwtConfig.JWTAccessSecret, i.jwtConfig.JWTRefreshSecret, i.jwtConfig.JWTAccessLifetime, i.jwtConfig.JWTRefreshLifetime)
	if err != nil {
		return dto.LoginResponse{Status: "JWT generation error"}, errors.New("JWT generation error")
	}

	err = i.repo.SaveRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(i.jwtConfig.JWTRefreshLifetime))
	if err != nil {
		log.Warn().Err(err).Send()
		return dto.LoginResponse{}, errors.New("failed to save refresh token")
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Status:       "ok",
	}, nil
}

func (i *UserUseCaseInteractor) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	_, err = jwt.VerifyToken(refreshToken, i.jwtConfig.JWTRefreshSecret)
	if err != nil {
		return "", "", err
	}

	dbRefreshToken, err := i.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) {
		i.repo.RevokeRefreshToken(ctx, refreshToken)
		return "", "", errors.New("refresh token has been revoked")
	}

	if !dbRefreshToken.RevokedAt.IsZero() {
		return "", "", errors.New("refresh token has been revoked")
	}

	user, err := i.repo.FindById(ctx, dbRefreshToken.UserID)
	if err != nil {
		return "", "", errors.New("associated user not found or inactive")
	}

	organizationID, err := i.repo.GetOrganizationID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			organizationID = uuid.Nil
		} else {
			return "", "", err
		}
	}

	err = i.repo.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	newAccessToken, newRefreshToken, err = jwt.GenerateTokens(
		dbRefreshToken.UserID,
		user.Role,
		organizationID,
		i.jwtConfig.JWTAccessSecret,
		i.jwtConfig.JWTRefreshSecret,
		i.jwtConfig.JWTAccessLifetime,
		i.jwtConfig.JWTRefreshLifetime,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}
	err = i.repo.SaveRefreshToken(ctx,
		dbRefreshToken.UserID,
		newRefreshToken,
		time.Now().Add(i.jwtConfig.JWTRefreshLifetime),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (i *UserUseCaseInteractor) VerifyUser(ctx context.Context, email string, code string) error {
	return nil
}
