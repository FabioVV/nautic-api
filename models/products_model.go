package models

import "time"

type CreateBoatRequest struct {
	Model   string `json:"Model" validate:"required"`
	NewUsed string `json:"NewUsed" validate:"required"`
}

type Boat struct {
	Id               int64     `json:"id"`
	Model            string    `json:"model"`
	PriceSell        *float64  `json:"selling_price"`
	Cost             *float64  `json:"cost"`
	Itens            *string   `json:"itens"`
	Hours            *string   `json:"hours"`
	Year             *string   `json:"year"`
	NewUsed          *string   `json:"new_used"`
	CabOpen          *string   `json:"cab_open"`
	Capacity         *string   `json:"capacity"`
	NightCapacity    *int      `json:"night_capacity"`
	Length           *float64  `json:"lenght"`
	Beam             *float64  `json:"beam"`
	Draft            *float64  `json:"draft"`
	Weight           *float64  `json:"weight"`
	Trim             *float64  `json:"trim"`
	FuelTankCapactiy *float64  `json:"fuel_tank_capacity"`
	Active           string    `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Accessory struct {
	Id              int64     `json:"id"`
	Model           string    `json:"model"`
	PriceBuy        float64   `json:"price_buy"`
	PriceSell       float64   `json:"price_sell"`
	Details         string    `json:"details"`
	Active          string    `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IdAccessoryType int32     `json:"AccessoryTypeId,omitempty"`
	Type            string    `json:"AccessoryType,omitempty"`
}

type Engine struct {
	Id         int64     `json:"id"`
	Model      string    `json:"model"`
	PriceSell  float64   `json:"selling_price"`
	Type       string    `json:"type"`
	Weight     float64   `json:"weight"`
	Rotation   string    `json:"rotation"`
	Power      float64   `json:"power"`
	Cylinders  float64   `json:"cylinders"`
	Command    string    `json:"command"`
	Clocks     float64   `json:"clocks"`
	Tempo      float64   `json:"tempo"`
	FuelType   string    `json:"fuel_type"`
	Active     string    `json:"active"`
	Propulsion *string   `json:"propulsion"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateAccessoryRequest struct {
	Model           *string  `json:"model,omitempty" validate:"required"`
	PriceBuy        *float64 `json:"PriceBuy,omitempty"`
	PriceSell       *float64 `json:"PriceSell,omitempty"`
	Details         *string  `json:"details,omitempty" validate:"required"`
	IdAccessoryType *int32   `json:"AccessoryTypeId,omitempty" validate:"required"`
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
	Model           string  `json:"Model" validate:"required"`
	PriceBuy        float64 `json:"PriceBuy,omitempty"`
	PriceSell       float64 `json:"PriceSell,omitempty"`
	Details         string  `json:"Details" validate:"required"`
	IdAccessoryType int32   `json:"AccessoryTypeId" validate:"required"`
}

type CreateEngineRequest struct {
	Model      string   `json:"Model" validate:"required"`
	PriceSell  *float64 `json:"SellingPrice" validate:"required"`
	Type       *string  `json:"Type" validate:"required"`
	Propulsion *string  `json:"Propulsion" validate:"required"`
	Weight     *float64 `json:"Weight"`
	Rotation   *string  `json:"Rotation"`
	Power      *float64 `json:"Power"`
	Cylinders  *float64 `json:"Cylinders"`
	Command    *string  `json:"Command" validate:"required"`
	Clocks     *float64 `json:"Clocks"`
	Tempo      *float64 `json:"Tempo" validate:"required"`
	FuelType   *string  `json:"FuelType" validate:"required"`
}
