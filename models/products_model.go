package models

import "time"

type Accessory struct {
	Id        int64     `json:"id"`
	Model     string    `json:"model"`
	PriceBuy  float64   `json:"price_buy"`
	PriceSell float64   `json:"price_sell"`
	Details   string    `json:"details"`
	Active    string    `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateAccessoryTypeRequest struct {
	Type *string `json:"type" validate:"required"`
}

type CreateAccessoryTypeRequest struct {
	Type string `json:"type" validate:"required"`
}

type AccessoryType struct {
	Id        int64     `json:"id"`
	Type      string    `json:"type"`
	Active    string    `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
