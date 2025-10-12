package entity

import (
	"github.com/google/uuid"
	"time"
)

type RolePermission struct {
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	CreatedAt    time.Time
}
