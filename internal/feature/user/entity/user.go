package entity

import "time"

type User struct {
	ID             uint64
	FullName       string
	Username       string
	HashedPassword string
	Email          string
	Phone          string
	Role           string
	CreatedAt      time.Time
}
