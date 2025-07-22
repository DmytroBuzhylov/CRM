package dto

type CreateOrganizationRequest struct {
	OwnerUserID uint64 `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateOrganizationResponse struct {
	Status string `json:"status"`
}
