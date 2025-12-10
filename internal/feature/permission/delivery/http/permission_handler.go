package http

import (
	"Test/internal/feature/permission/usecase"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	uc usecase.PermissionUseCase
}

func NewPermissionHandler(uc usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{uc: uc}
}

func (h *PermissionHandler) CreateOrganizationRoleHandler(c *gin.Context) {}

func (h *PermissionHandler) UpdatePermissionToOrganizationRoleHandler(c *gin.Context) {}

func (h *PermissionHandler) AdmitRoleToStaffHandler(c *gin.Context) {}
