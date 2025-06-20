package http

import (
	"Test/internal/feature/user/interface_adapters/dto"
	"Test/internal/feature/user/usecase"
	"Test/internal/pkg/jwt"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserHandler struct {
	createUserUC usecase.CreateUser
	getUserUC    usecase.GetUser
	updateUserUC usecase.UpdateUser
	deleteUserUC usecase.DeleteUser
}

func NewUserHandler(
	createUserUC usecase.CreateUser,
	getUserUC usecase.GetUser,
	updateUserUC usecase.UpdateUser,
	deleteUserUC usecase.DeleteUser,
) *UserHandler {
	return &UserHandler{
		createUserUC: createUserUC,
		getUserUC:    getUserUC,
		updateUserUC: updateUserUC,
		deleteUserUC: deleteUserUC,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.createUserUC.Create(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cookie error",
		})
		return
	}

}
