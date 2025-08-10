package http

import (
	"Test/internal/feature/inventory/interface_adapter/dto"
	"Test/internal/feature/inventory/usecase"
	"Test/internal/pkg/storage"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type InventoryHandler struct {
	inventoryUseCase usecase.InventoryUseCase
}

func NewInventoryHandler(inventoryUseCase usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{inventoryUseCase: inventoryUseCase}
}

func (h *InventoryHandler) CreateIngredientHandler(c *gin.Context) {
	var (
		req dto.CreateIngredientRequest
	)

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The user has not joined the organization",
		})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to bind request: " + err.Error(),
		})
		return
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file: " + err.Error(),
		})
		return
	}
	defer file.Close()
	if fileHeader.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File is too large. Max size is 2MB",
		})
		return
	}

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read file buffer",
		})
		return
	}

	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported file type. Only JPEG, PNG, and GIF are allowed",
		})
		return

	}

	fileOpts := storage.StorageFileOptions{
		File:        file,
		FileSize:    fileHeader.Size,
		ContentType: contentType,
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.CreateIngredient(ctx, req, fileOpts)
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

func (h *InventoryHandler) GetOrganizationIngredients(c *gin.Context) {
	var req dto.GetAllIngredientsRequest

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The user has not joined the organization",
		})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.GetAllIngredients(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *InventoryHandler) GetIngredientHandler(c *gin.Context) {
	var req dto.GetIngredientRequest
	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "The user has not joined the organization",
		})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ingredientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ingredient id parse error",
		})
		return
	}
	req.ID = ingredientID

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.GetIngredient(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *InventoryHandler) UpdateIngredientHandler(c *gin.Context) {
	var req dto.UpdateIngredientRequest

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "The user has not joined the organization",
		})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ingredientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Ingredient id parse error",
		})
		return
	}
	req.ID = ingredientID

	if err = c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.UpdateIngredient(ctx, req)
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

func (h *InventoryHandler) DeleteIngredientHandler(c *gin.Context) {
	var (
		req dto.DeleteIngredientRequest
		err error
	)

	ingredientID := c.Param("id")
	if ingredientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	req.ID, err = uuid.Parse(ingredientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	organizationID, ok := c.Get("organization_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "The user has not joined the organization",
		})
		return
	}
	req.OrganizationID = organizationID.(uuid.UUID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.DeleteIngredient(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": resp.Status,
	})
}

func (h *InventoryHandler) DeleteManyIngredientsHandler(c *gin.Context) {
	var req dto.DeleteManyIngredientsRequest
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.inventoryUseCase.DeleteManyIngredients(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": resp.Status,
	})
}
