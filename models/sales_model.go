package models

import "time"

type Customer struct {
	Id             int64   `json:"id"`
	UserId         int64   `json:"id_user"`
	MeanComId      int32   `json:"id_mean_communication"`
	SellerName     string  `json:"seller_name"`
	MeamComName    string  `json:"com_name"`
	Name           string  `json:"customer_name"`
	Email          *string `json:"customer_email"`
	Phone          *string `json:"customer_phone"`
	BirthDate      *string `json:"birthdate"`
	PfPj           *string `json:"pf_pj"`
	Cpf            *string `json:"cpf"`
	Cnpj           *string `json:"cnpj"`
	Cep            *string `json:"cep"`
	Street         *string `json:"street"`
	Neighborhood   *string `json:"neighborhood"`
	City           *string `json:"city"`
	State          *string `json:"state"`
	Complement     *string `json:"complement"`
	Qualified      *string `json:"qualified"`
	QualifiedType  *string `json:"qualified_type,omitempty"`
	Active         *string `json:"active"`
	ActiveContact  *string `json:"active_contact"`
	SuspectOfFraud *string `json:"suspect_of_fraud"`

	//Negotiation data, sometimes you might not need this and in others it might be null
	EstimatedValue *float64 `json:"estimated_value"`
	NewUsed        *string  `json:"new_used,omitempty"`
	CustomerCity   *string  `json:"customer_city"`
	NavigationCity *string  `json:"customer_nav_city,omitempty"`
	BoatCapacity   *int16   `json:"boat_cap_needed,omitempty"`
	CabinatedOpen  *string  `json:"cabinated_open,omitempty"`
	MaxPesBoat     *string  `json:"boat_length_max"`
	MinPesBoat     *string  `json:"boat_length_min"`
}

type CustomerRequest struct {
	Name         *string    `json:"Name" validate:"required"`
	Email        *string    `json:"Email" validate:"required"`
	Phone        *string    `json:"Phone" validate:"required"`
	BirthDay     *time.Time `json:"Birthday"`
	PfPj         *string    `json:"PfPj" validate:"required"`
	Cpf          *string    `json:"Cpf" `
	Cnpj         *string    `json:"Cnpj"`
	Cep          *string    `json:"Cep" validate:"required"`
	Street       *string    `json:"Street" validate:"required"`
	Neighborhood *string    `json:"Neighborhood" validate:"required"`
	City         *string    `json:"City" validate:"required"`
	Complement   *string    `json:"Complement" validate:"required"`
	State        *string    `json:"State" validate:"required"`

	HasBoat   *string `json:"HasBoat,omitempty"`
	WhichBoat *string `json:"WhichBoat,omitempty"`
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

	HasBoat      *string `json:"has_boat"`
	HasBoatWhich *string `json:"has_boat_which"`
	MaxPesBoat   *string `json:"boat_length_max"`
	MinPesBoat   *string `json:"boat_length_min"`

	Stage         int64  `json:"stage"`
	Qualified     string `json:"qualified"`
	QualifiedType string `json:"qualified_type"`

	HasPassed24Hrs *bool `json:"has_passed_24hrs"`
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
	BusinessId   *int64    `json:"id_business"`
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

	HasBoat    *string `json:"HasBoat,omitempty"`
	WhichBoat  *string `json:"WhichBoat,omitempty"`
	MinPesBoat *int    `json:"MinPesBoat,omitempty"`
	MaxPesBoat *int    `json:"MaxPesBoat,omitempty"`

	UserId *int64 `json:"UserId" validate:"required"`
}

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
