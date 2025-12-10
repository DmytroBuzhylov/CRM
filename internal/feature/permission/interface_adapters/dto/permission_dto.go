package dto

import (
	"github.com/google/uuid"
)

type AdmitRoleToStaffRequestDTO struct {
	StaffID        uuid.UUID
	RoleIDsToAdmit []uuid.UUID
}
