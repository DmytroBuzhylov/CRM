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

func NewUser(fullName string, username string, hashedPassword string, email string, phone string, role string) *User {
	now := time.Now()
	return &User{
		FullName:       fullName,
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
		Phone:          phone,
		Role:           role,
		CreatedAt:      now,
	}
}
