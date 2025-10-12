package http

import (
	"Test/internal/feature/notification/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type NotificationHTTPHandler struct {
	uc usecase.NotificationUseCase
}

func NewNotificationHTTPHandler(uc usecase.NotificationUseCase) *NotificationHTTPHandler {
	return &NotificationHTTPHandler{uc: uc}
}

func (h *NotificationHTTPHandler) GetUserNotifications(c *gin.Context) {
	userIdStr, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "absent user_id",
		})
		return
	}
	userID := userIdStr.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	resp, err := h.uc.GetUserNotifications(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
