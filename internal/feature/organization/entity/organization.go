package entity

import (
	"github.com/google/uuid"
	"time"
)

type Organization struct {
	ID          uuid.UUID
	Name        string
	Description string
	OwnerUserID uuid.UUID
	CreatedAt   time.Time
}

func NewOrganization(name string, description string, ownerUserID uuid.UUID) *Organization {
	now := time.Now()
	return &Organization{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerUserID: ownerUserID,
		CreatedAt:   now,
	}
}
