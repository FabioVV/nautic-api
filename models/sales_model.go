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

type NegotiationHistory struct {
	Id           int64     `json:"id"`
	ComMeanId    *int64    `json:"id_mean_communication"`
	MeamComName  string    `json:"com_name"`
	UserId       *int64    `json:"id_user"`
	CustomerId   *int64    `json:"id_customer"`
	CustomerName string    `json:"customer_name"`
	Description  *string   `json:"description"`
	Stage        *int64    `json:"stage"`
	DateCreated  time.Time `json:"created_at"`
}

type CreateNegotiationRequest struct {
	Name           *string  `json:"Name"`
	Email          *string  `json:"Email"`
	Phone          *string  `json:"Phone"`
	EstimatedValue *float64 `json:"EstimatedValue" validate:"required"`
	BoatName       *string  `json:"BoatName" validate:"required"`
	NewUsed        *string  `json:"NewUsed,omitempty"`
	Qualified      *string  `json:"Qualified,omitempty"`
	QualifiedType  *string  `json:"QualifiedType,omitempty"`
	City           *string  `json:"City,omitempty"`
	NavigationCity *string  `json:"NavigationCity,omitempty"`
	BoatCapacity   *int16   `json:"BoatCapacity,omitempty"`
	CabinatedOpen  *string  `json:"CabinatedOpen,omitempty"`
	ComMeanId      *int32   `json:"ComMeanId"`
	UserId         *int64   `json:"UserId" validate:"required"`
}

// acoForm = this.formBuilder.group({
//     Description: ['', [Validators.required]],
//     ComMeanName: ['', [Validators.required]],
//     ComMeanId: ['', [Validators.required]],
//     UserId: ['', []],
// })

type CreateNegotiationHistoryRequest struct {
	Description *string `json:"Description" validate:"required"`
	ComMeanId   *int64  `json:"ComMeanId" validate:"required"`
	UserId      *int64  `json:"UserId" validate:"required"`
	CustomerId  *int64  `json:"CustomerId" validate:"required"`
	Stage       *int64  `json:"Stage" validate:"required"`
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
