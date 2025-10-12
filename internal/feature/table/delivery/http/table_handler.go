package http

import (
	"Test/internal/feature/table/interface_adapters/dto"
	"Test/internal/feature/table/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type TableHandler struct {
	uc usecase.TableUseCase
}

func NewTableHandler(uc usecase.TableUseCase) *TableHandler {
	return &TableHandler{uc: uc}
}

func (h *TableHandler) CreateTableHandler(c *gin.Context) {
	var req dto.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSONP(http.StatusUnauthorized, gin.H{})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := h.uc.CreateTable(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}
