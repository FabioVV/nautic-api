package models

import "time"

type SalesOrder struct {
	Id         *int64  `json:"id"`
	CustomerId *int64  `json:"id_customer"`
	UserId     *int64  `json:"id_user"` // seller
	Uuid       *string `json:"uuid"`

	StatusType *string `json:"status_type"`

	CustomerName  *string `json:"customer_name"`
	CustomerEmail *string `json:"customer_email"`
	CustomerPhone *string `json:"customer_phone"`

	OrderBoatModel   *string `json:"OrderBoatModel"`
	OrderBoatId      *string `json:"OrderBoatId"`
	OrderEngineModel *string `json:"OrderEngineModel"`
	OrderEngineId    *string `json:"OrderEngineId"`

	OrderBoatPrice   *float64 `json:"OrderBoatPrice"`
	OrderEnginePrice *float64 `json:"OrderEnginePrice"`

	TotalItensPrice *float64 `json:"TotalItensPrice"`

	SellerName *string `json:"seller_name"`

	PfPj         *string `json:"PfPj" validate:"required"`
	Cpf          *string `json:"Cpf" `
	Cnpj         *string `json:"Cnpj"`
	Cep          *string `json:"Cep" validate:"required"`
	Street       *string `json:"Street" validate:"required"`
	Neighborhood *string `json:"Neighborhood" validate:"required"`
	City         *string `json:"City" validate:"required"`
	Complement   *string `json:"Complement" validate:"required"`
	State        *string `json:"State" validate:"required"`

	Status           *string    `json:"status"`
	Done             *string    `json:"done"`
	DiscountedAmount *float64   `json:"discounted_amount" `
	TotalValue       *float64   `json:"total_value"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	DeliveryAt       *time.Time `json:"delivery_at"`
}

type SalesOrderItem struct {
	Id *int64 `json:"id"`

	AccessoryId *int64  `json:"id_accessory"`
	Model       *string `json:"model"`

	UnitPrice       *float64 `json:"unit_price"`
	Quantity        *int64   `json:"qty"`
	Discount        *float64 `json:"discount"`
	TotalPriceItens *float64 `json:"total_price_itens"`
}

type UpdateSalesOrderItemQty struct {
	Qty *int64 `json:"qty"`
}

type SalesOrderFile struct {
	Id   int64  `json:"id"`
	Path string `json:"path"`
	Type *int   `json:"type"`
}

type UpdateSalesOrderFileType struct {
	Type *int `json:"type" validate:"required"`
}

type LostNegotiationReport struct {
	Id                     *int64     `json:"id"`
	BusinessId             *int64     `json:"id_business"`
	CustomerName           *string    `json:"customer_name"`
	Motive                 *string    `json:"motive"`
	WhichBoat              *string    `json:"which_boat"`
	OurBoatOffered         *string    `json:"our_boat_offered"`
	CustomerGotAnotherBoat *string    `json:"customer_got_another_boat"`
	CreatedAt              *time.Time `json:"created_at"`
}
