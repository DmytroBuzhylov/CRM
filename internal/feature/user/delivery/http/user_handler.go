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
	createUserUC usecase.CreateUser
	getUserUC    usecase.GetUser
	updateUserUC usecase.UpdateUser
	deleteUserUC usecase.DeleteUser
	jwtConfig    config.JWTConfig
}

func NewUserHandler(
	createUserUC usecase.CreateUser,
	getUserUC usecase.GetUser,
	//updateUserUC usecase.UpdateUser,
	//deleteUserUC usecase.DeleteUser,
	jwtConfig config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		//updateUserUC: updateUserUC,
		//deleteUserUC: deleteUserUC,
		jwtConfig: jwtConfig,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		log.Warn().Err(err)
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.createUserUC.Create(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		log.Warn().Err(err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unauthorized: Refresh token missing",
		})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	newAccessToken, newRefreshToken, err := h.getUserUC.Refresh(ctx, refreshTokenString)
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

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.getUserUC.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
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
