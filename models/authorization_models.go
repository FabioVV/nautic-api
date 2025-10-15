package models

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type Role struct {
	Id        *int64     `json:"id"`
	Name      *string    `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UserPermission struct {
	PermissionId     *int64 `json:"id_permission"`
	UserPermissionID *int64 `json:"id_role"`

	Module      *string `json:"module"`
	Description *string `json:"description"`

	HasPermission *bool `json:"has_permission"`
}

type RolePermission struct {
	PermissionId *int64 `json:"id_permission"`
	RoleID       *int64 `json:"id_role"`

	Module      *string `json:"module"`
	Description *string `json:"description"`
	RoleName    *string `json:"name"`

	HasPermission *bool `json:"has_permission"`
}

type Permission struct {
	Id          *int64  `json:"id"`
	Module      *string `json:"module"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	Url         *string `json:"url"`
}

type CreateRole struct {
	Name *string `json:"Name" validate:"required"`
}
