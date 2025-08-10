package http

import (
	"Test/internal/feature/organization/interface_adapters/dto"
	"Test/internal/feature/organization/usecase"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type OrganizationHandler struct {
	organizationUC usecase.OrganizationUseCase
	invitationUC   usecase.InvitationUseCase
}

func NewOrganizationHandler(organizationUC usecase.OrganizationUseCase, invitationUC usecase.InvitationUseCase) *OrganizationHandler {
	return &OrganizationHandler{organizationUC: organizationUC, invitationUC: invitationUC}
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
	req.OwnerUserID = userID.(uuid.UUID)

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

func (h *OrganizationHandler) CreateInvite(c *gin.Context) {
	var req dto.CreateInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
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

	resp, err := h.invitationUC.GenerateInvitation(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": resp.Code,
	})
}

func (h *OrganizationHandler) AcceptInvite(c *gin.Context) {
	var req dto.AcceptInvitationRequest
	req.Code = c.Param("code")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	req.UserID = userID.(uuid.UUID)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	resp, err := h.invitationUC.AcceptInvitation(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		log.Info().Err(err).Send()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": resp.Status,
	})
}
