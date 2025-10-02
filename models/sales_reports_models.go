package models

import "time"

type NegotiationReport struct {
	Id             *int64   `json:"id"`
	CustomerId     *int64   `json:"id_customer"`
	MeanComId      *int32   `json:"id_mean_communication"`
	MeamComName    *string  `json:"com_name"`
	Name           *string  `json:"customer_name"`
	Email          *string  `json:"customer_email"`
	Phone          *string  `json:"customer_phone"`
	BoatName       *string  `json:"boat_name"`
	EstimatedValue *float64 `json:"estimated_value"`
	Stage          *int     `json:"stage"`

	DaysSinceStageChange *int       `json:"days_since_stage_change"`
	LastHistoryAt        *time.Time `json:"last_history_at"`
	DaysSinceLastHistory *int       `json:"days_since_last_history"`
}

type SalesOrdersReport struct {
	Id                *int64 `json:"id"`
	CustomerId        *int64 `json:"id_customer"`
	UserId            *int64 `json:"id_user"`
	BusinessHistoryId *int64 `json:"id_business_history"`

	CustomerName *string `json:"customer_name"`
	SellerName   *string `json:"seller_name"`

	Status           *string    `json:"status"`
	Done             *string    `json:"done"`
	DiscountedAmount *float64   `json:"discounted_amount" `
	TotalValue       *float64   `json:"total_value"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeliveryAt       *time.Time `json:"delivery_at"`
}
