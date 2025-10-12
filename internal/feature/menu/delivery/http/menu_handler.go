package http

import (
	"Test/internal/feature/menu/interface_adapter/dto"
	"Test/internal/feature/menu/usecase"
	"Test/pkg/utils"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type MenuHandler struct {
	uc usecase.MenuItemUseCase
}

func NewMenuHandler(uc usecase.MenuItemUseCase) *MenuHandler {
	return &MenuHandler{uc: uc}
}

func (h *MenuHandler) CreateMenuItemHandler(c *gin.Context) {

	jsonData := c.PostForm("data")
	if jsonData == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing JSON data in form field 'data'"})
		return
	}

	var req *dto.MenuItemWithRecipeDTO
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MenuItem == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "menu_item is required"})
		return
	}

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the user does not belong to the organization",
		})
		return
	}
	req.MenuItem.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fileOpts, err := utils.FileVerifier(file, fileHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.uc.CreateMenuItem(ctx, req, fileOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Status": "ok",
	})
}

func (h *MenuHandler) UpdateMenuItemHandler(c *gin.Context) {
	var (
		req *dto.MenuItemWithRecipeDTO
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.uc.UpdateMenuItem(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": resp.Status,
	})
}
