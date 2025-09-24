package models

import "time"

type Customer struct {
	Id            int64   `json:"id"`
	UserId        int64   `json:"id_user"`
	MeanComId     int32   `json:"id_mean_communication"`
	SellerName    string  `json:"seller_name"`
	MeamComName   string  `json:"com_name"`
	Name          string  `json:"customer_name"`
	Email         *string `json:"customer_email"`
	Phone         *string `json:"customer_phone"`
	BirthDate     *string `json:"birthdate"`
	PfPj          *string `json:"pf_pj"`
	Cpf           *string `json:"cpf"`
	Cnpj          *string `json:"cnpj"`
	Cep           *string `json:"cep"`
	Street        *string `json:"street"`
	Neighborhood  *string `json:"neighborhood"`
	City          *string `json:"city"`
	Complement    *string `json:"complement"`
	Qualified     *string `json:"qualified"`
	Active        *string `json:"active"`
	ActiveContact *string `json:"active_contact"`
}

type Negotiation struct {
	Id                 int64    `json:"id"`
	CustomerId         int64    `json:"id_customer"`
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
	QualifiedType      string   `json:"qualified_type"`
}

type CreateNegotiationRequest struct {
	Name           *string  `json:"Name" validate:"required"`
	Email          *string  `json:"Email" validate:"required"`
	Phone          *string  `json:"Phone" validate:"required"`
	EstimatedValue *float64 `json:"EstimatedValue" validate:"required"`
	BoatName       *string  `json:"BoatName" validate:"required"`
	Qualified      *string  `json:"Qualified,omitempty"`
	QualifiedType  *string  `json:"QualifiedType,omitempty"`
	ComMeanId      *int32   `json:"ComMeanId"`
	UserId         *int64   `json:"UserId" validate:"required"`
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
