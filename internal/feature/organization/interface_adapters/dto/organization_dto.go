package dto

type CreateOrganizationRequest struct {
	OwnerUserID uint64
	Name        string
	Description string
}

type CreateOrganizationResponse struct {
	Status string
}
