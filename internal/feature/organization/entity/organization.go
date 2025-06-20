package entity

import "time"

type Organization struct {
	ID          uint64
	Name        string
	Description string
	OwnerUserID uint64
	CreatedAt   time.Time
}

func NewOrganization(name string, description string, ownerUserID uint64) *Organization {
	now := time.Now()
	return &Organization{
		Name:        name,
		Description: description,
		OwnerUserID: 0,
		CreatedAt:   now,
	}
}
