package models

import "time"

type Role struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permission struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Active string `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
