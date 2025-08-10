package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=3"`
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,startswith=+380,len=13"`
}

type CreateUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string
	Status       string `json:"status"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Method   string `json:"method" binding:"required"` // email, username, phone
}

type GetUserRequest struct {
	ID       uuid.UUID `json:"-"` // handler receives data from jwt
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type GetUserResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Role     string    `json:"role"`
}

type UpdateUserRequest struct {
	ID       uuid.UUID `json:"-"` // handler receives data from jwt
	FullName string    `json:"full_name"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
}

type UpdateUserResponse struct {
	Status bool `json:"status"`
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"-"` // handler receives data from jwt
}

type DeleteUserResponse struct {
	Status bool `json:"status"`
}
