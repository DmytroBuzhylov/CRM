package usecase

import (
	"Test/config"
	"Test/internal/feature/user/entity"
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/repository"
	"Test/internal/feature/user/repository/postgres"
	"Test/internal/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"time"
)

type getUserInteractor struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewGetUserInteractor(userRepo repository.UserRepository) *getUserInteractor {
	return &getUserInteractor{userRepo: userRepo}
}

func (i *getUserInteractor) GetById(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *getUserInteractor) GetByUsername(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *getUserInteractor) GetByEmail(ctx context.Context, req dto.GetUserRequest) (dto.GetUserResponse, error) {
	return dto.GetUserResponse{}, fmt.Errorf("")
}

func (i *getUserInteractor) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	var (
		user entity.User
		err  error
	)
	switch req.Method {
	case "email":
		user, err = i.userRepo.FindByEmail(ctx, req.Email)

	case "phone":
		user, err = i.userRepo.FindByPhone(ctx, req.Phone)

	case "username":
		user, err = i.userRepo.FindByUsername(ctx, req.Username)

	default:
		return dto.LoginResponse{Status: "incorrect method"}, fmt.Errorf("incorrect method")
	}

	if err != nil {
		return dto.LoginResponse{Status: "db error"}, err
	}
	if !postgres.CheckPasswordHash(req.Password, user.HashedPassword) {
		return dto.LoginResponse{Status: "incorrect password"}, err
	}

	accessToken, refreshToken, err := jwt.GenerateTokens(user.ID, user.Role, i.jwtConfig.JWTAccessSecret, i.jwtConfig.JWTRefreshSecret, i.jwtConfig.JWTAccessLifetime, i.jwtConfig.JWTRefreshLifetime)
	if err != nil {
		return dto.LoginResponse{Status: "JWT generation error"}, err
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Status:       "ok",
	}, nil
}

func (i *getUserInteractor) Refresh(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	_, err = jwt.VerifyToken(refreshToken, i.jwtConfig.JWTRefreshSecret)
	if err != nil {
		return "", "", err
	}

	dbRefreshToken, err := i.userRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if time.Now().After(dbRefreshToken.ExpiresAt) {
		i.userRepo.RevokeRefreshToken(ctx, refreshToken)
		return "", "", errors.New("refresh token has been revoked")
	}

	if !dbRefreshToken.RevokedAt.IsZero() {
		return "", "", errors.New("refresh token has been revoked")
	}

	user, err := i.userRepo.FindById(ctx, dbRefreshToken.UserID)
	if err != nil {
		return "", "", errors.New("associated user not found or inactive")
	}

	err = i.userRepo.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	newAccessToken, newRefreshToken, err = jwt.GenerateTokens(
		dbRefreshToken.UserID,
		user.Role,
		i.jwtConfig.JWTAccessSecret,
		i.jwtConfig.JWTRefreshSecret,
		i.jwtConfig.JWTAccessLifetime,
		i.jwtConfig.JWTRefreshLifetime,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %w", err)
	}
	err = i.userRepo.SaveRefreshToken(ctx,
		dbRefreshToken.UserID,
		newRefreshToken,
		time.Now().Add(i.jwtConfig.JWTRefreshLifetime),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}
