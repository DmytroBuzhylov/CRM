package http

import (
	"Test/internal/feature/order/interface_adapters/dto"
	"Test/internal/feature/order/usecase"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type OrderHandler struct {
	uc usecase.OrderUseCase
}

func NewOrderHandler(uc usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) CreateOrderHandler(c *gin.Context) {
	var req *dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSONP(http.StatusUnauthorized, gin.H{})
		return
	}
	req.Order.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err := h.uc.CreateOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func (h *OrderHandler) GetOrderHandler(c *gin.Context) {
	var (
		err error
		req *dto.GetOrderRequest
	)

	idString := c.Param("id")
	req.ID, err = uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.uc.GetOrder(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrganizationsOrdersHandler(c *gin.Context) {
	var (
		req = &dto.GetOrganizationOrdersRequest{}
	)

	organizationID, ok := c.Get("organization_id")
	if !ok || organizationID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid organization id",
		})
		return
	}

	log.Info().Msg(fmt.Sprintf("%v", organizationID))
	req.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.uc.GetOrganizationOrders(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
