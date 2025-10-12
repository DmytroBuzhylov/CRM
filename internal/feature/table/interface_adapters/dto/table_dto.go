package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateTableRequest struct {
	OrganizationID uuid.UUID `json:"-"`
	Number         uint64    `json:"number" binding:"required"`
	Seats          uint64    `json:"seats" binding:"required"`
}

type GetTableRequest struct {
	ID             uuid.UUID `json:"ID" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
}

type GetOrganizationTablesRequest struct {
	OrganizationID uuid.UUID `json:"-"`
}

type GetTableResponse struct {
	ID             uuid.UUID `json:"ID"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Number         uint64    `json:"number"`
	Status         string    `json:"status"`
	Seats          uint64    `json:"seats"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateStatusRequest struct {
	ID             uuid.UUID `json:"ID" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
	Status         string    `json:"status"`
}

type UpdateSeatsRequest struct {
	ID             uuid.UUID `json:"ID" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
	Seats          uint64    `json:"seats"`
}

type UpdateNumberRequest struct {
	ID             uuid.UUID `json:"ID" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
	Number         uint64    `json:"number"`
}

type DeleteRequest struct {
	ID             uuid.UUID `json:"ID" binding:"required"`
	OrganizationID uuid.UUID `json:"-"`
}
