package entity

import "github.com/google/uuid"

type OrganizationUser struct {
	OrganizationID uuid.UUID
	UserID         uint64
	Role           string
}
