package http

import (
	"Test/internal/feature/staff/usecase"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	uc usecase.StaffUseCase
}

func NewStaffHandler(uc usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{uc: uc}
}

func (h *StaffHandler) CreateInviteHandler(c *gin.Context) {
	
}

func (h *StaffHandler) AcceptInviteHandler(c *gin.Context) {}

func (h *StaffHandler) DismissStaffHandler(c *gin.Context) {}
