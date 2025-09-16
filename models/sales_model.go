package models

import "time"

// {
//     "Name": "assasda",
//     "Email": "sadsdasdasda",
//     "Phone": "12-31231-2313",
//     "EstimatedValue": 11111.02,
//     "Qualified": "N",
//     "QualifiedType": "",
//     "BoatName": "asdsaddasdsadsa",
//     "ComMeanName": "Facebook #updated",
//     "ComMeanId": 2
// }

type CreateNegotiationRequest struct {
	Name           *string  `json:"Name,omitempty" validate:"required"`
	Email          *string  `json:"Email,omitempty" validate:"required"`
	Phone          *string  `json:"Phone,omitempty" validate:"required"`
	EstimatedValue *float64 `json:"EstimatedValue,omitempty" validate:"required"`
	BoatName       *string  `json:"BoatName,omitempty" validate:"required"`
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
