package http

import (
	"Test/config"
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/usecase"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type AuthHandler struct {
	uc        usecase.UserUseCase
	jwtConfig config.JWTConfig
}

func NewUserHandler(
	uc usecase.UserUseCase,
	jwtConfig config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		uc:        uc,
		jwtConfig: jwtConfig,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		log.Warn().Err(err).Send()
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.uc.Create(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Warn().Err(err).Send()
		return
	}

	log.Info().Msg("user has been registered")

	c.SetCookie(
		"refresh_token",
		resp.RefreshToken,
		int(h.jwtConfig.JWTRefreshLifetime.Seconds()),
		"/api/v1/auth/refresh",
		"",
		false,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"access_token": resp.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.jwtConfig.JWTAccessLifetime.Seconds()),
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: Refresh token missing",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	newAccessToken, newRefreshToken, err := h.uc.Refresh(ctx, refreshTokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unauthorized: %v", err.Error())})
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		int(h.jwtConfig.JWTRefreshLifetime.Seconds()),
		"/api/v1/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.jwtConfig.JWTAccessLifetime.Seconds()),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}
	if req.Method != "email" && req.Method != "phone" && req.Method != "username" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid method",
		})
		return
	}
	if req.Password == "" || len(req.Password) <= 7 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid password",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.uc.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		resp.RefreshToken,
		int(h.jwtConfig.JWTRefreshLifetime.Seconds()),
		"/api/v1/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": resp.AccessToken,
		"token_type":   "Bearer",
		"expires_in":   int(h.jwtConfig.JWTAccessLifetime.Seconds()),
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	//defer cancel()
	//
	//err := h.getUserUC.

}

func (h *AuthHandler) Login2FA(c *gin.Context) {

}
