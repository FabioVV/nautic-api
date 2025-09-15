package models

import "time"

type CreateCommunicationMeanRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCommunicationMeaneRequest struct {
	Name *string `json:"name" validate:"required"`
}

type CommunicationMean struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	Active       string    `json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
