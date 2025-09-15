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
	IdAccessoryType int32  `json:"AccessoryTypeId,omitempty"`
	Type string  `json:"AccessoryType,omitempty"`
}

type UpdateAccessoryRequest struct {
	Model     *string    `json:"model,omitempty" validate:"required"`
	PriceBuy  *float64   `json:"PriceBuy,omitempty"`
	PriceSell *float64   `json:"PriceSell,omitempty"`
	Details   *string    `json:"details,omitempty" validate:"required"`
	IdAccessoryType *int32  `json:"AccessoryTypeId,omitempty" validate:"required"`
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

type CreateAccessoryRequest struct {
	Model     string    `json:"Model" validate:"required"`
	PriceBuy  float64   `json:"PriceBuy,omitempty"`
	PriceSell float64   `json:"PriceSell,omitempty"`
	Details   string    `json:"Details" validate:"required"`
	IdAccessoryType int32  `json:"AccessoryTypeId" validate:"required"`
}
