package entity

import (
	"github.com/google/uuid"
	"github.com/maniartech/signals"
	"time"
)

var (
	UserCreated       = signals.New[User]()
	UserRegistered    = signals.New[User]()
	UserTryingToLogin = signals.New[User]()
	UserLoggedIn      = signals.New[User]()
	UserUpdated       = signals.New[User]()
)

type User struct {
	ID             uuid.UUID
	FullName       string
	Username       string
	HashedPassword string
	Email          string
	Phone          string
	Role           string
	CreatedAt      time.Time
}

func NewUser(fullName string, username string, hashedPassword string, email string, phone string, role string) User {
	now := time.Now()
	return User{
		ID:             uuid.New(),
		FullName:       fullName,
		Username:       username,
		HashedPassword: hashedPassword,
		Email:          email,
		Phone:          phone,
		Role:           role,
		CreatedAt:      now,
	}
}
