package entity

import "time"

type Ingredient struct {
	ID              uint64
	OrganizationID  uint64
	Name            string
	Unit            string
	Quantity        uint64
	MinimumQuantity uint64
	ImageURL        string
	CreatedAt       *time.Time
}
