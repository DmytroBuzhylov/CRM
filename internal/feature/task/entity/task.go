package entity

import "time"

type Task struct {
	ID             uint64
	Name           string
	Description    string
	Priority       uint
	Status         string
	Deadline       *time.Time
	AssigneeID     *uint64
	ClientID       *uint64
	OrganizationID uint64
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}
