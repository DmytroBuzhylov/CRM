package http

import (
	"Test/internal/feature/organization/interface_adapters/dto"
	"Test/internal/feature/organization/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type OrganizationHandler struct {
	organizationUC usecase.OrganizationUsecase
}

func NewOrganizationHandler(organizationUC usecase.OrganizationUsecase) *OrganizationHandler {
	return &OrganizationHandler{organizationUC: organizationUC}
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{
			"error": "invalid body",
		})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSONP(http.StatusUnauthorized, gin.H{})
		return
	}
	req.OwnerUserID = userID.(uint64)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.organizationUC.Create(ctx, req)
	if err != nil {
		log.Err(err).Send()
		c.JSONP(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
