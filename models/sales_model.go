package models

import "time"

// SELECT id,
// 		id_customer,
//  		id_mean_communication,
// 		num_bussiness_customer,
// 		boat_name,
// 		estimated_value,
// 		max_estimated_value,
// 		customer_city,
// 		customer_navigation_city,
// 		boat_capacity_needed,
// 		new_used,
// 		cab_open,
// 		stage,
// 		qualified

type Negotiation struct {
	Id                 int64    `json:"id"`
	CustomerId         int32    `json:"id_customer"`
	MeanComId          int32    `json:"id_mean_communication"`
	Name               string   `json:"customer_name"`
	Email              string   `json:"customer_email"`
	Phone              string   `json:"customer_phone"`
	MeamComName        string   `json:"com_name"`
	BoatName           string   `json:"boat_name"`
	EstimatedValue     float64  `json:"estimated_value"`
	MaxEstimatedValue  *float64 `json:"max_estimated_value"`
	City               *string  `json:"customer_city"`
	NavigationCity     *string  `json:"customer_nav_city"`
	BoatCapacityNeeded *int32   `json:"boat_cap_needed"`
	NewUsed            *string  `json:"new_used"`
	CabOpen            *string  `json:"cab_open"`
	Stage              int64    `json:"stage"`
	Qualified          string   `json:"qualified"`
}

type CreateNegotiationRequest struct {
	Name           *string  `json:"Name" validate:"required"`
	Email          *string  `json:"Email" validate:"required"`
	Phone          *string  `json:"Phone" validate:"required"`
	EstimatedValue *float64 `json:"EstimatedValue" validate:"required"`
	BoatName       *string  `json:"BoatName" validate:"required"`
	Qualified      *string  `json:"Qualified,omitempty"`
	QualifiedType  *string  `json:"QualifiedType,omitempty"`
}

type CreateCommunicationMeanRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCommunicationMeaneRequest struct {
	Name *string `json:"name" validate:"required"`
}

type CommunicationMean struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Active    string    `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
