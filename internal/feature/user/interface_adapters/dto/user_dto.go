package dto

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type CreateUserResponse struct {
	ID   uint64 `json:"id"`
	Role string `json:"role"`
}

type GetUserRequest struct {
	ID       uint64 `json:"-"` // handler receives data from jwt
	Email    string `json:"email"`
	Username string `json:"username"`
}

type GetUserResponse struct {
	ID       uint64 `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	ID       uint64 `json:"-"` // handler receives data from jwt
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UpdateUserResponse struct {
	Status bool `json:"status"`
}

type DeleteUserRequest struct {
	ID uint64 `json:"-"` // handler receives data from jwt
}

type DeleteUserResponse struct {
	Status bool `json:"status"`
}
