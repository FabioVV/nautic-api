package models

import "time"

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:""`
	Email    *string `json:"email,omitempty" validate:"email"`
	Phone    *string `json:"phone,omitempty"`
	Active   *string `json:"active,omitempty"`
}

type UpdateUserWithPasswordRequest struct {
	Name     *string `json:"name,omitempty" validate:""`
	Email    *string `json:"email,omitempty" validate:"email"`
	Phone    *string `json:"phone,omitempty"`
	Active   *string `json:"active,omitempty"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
	OldPassword string `json:"old_password" validate:"required,min=8"`
}

type User struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Active       string    `json:"active"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
