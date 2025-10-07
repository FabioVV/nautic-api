package repositories

import (
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetSalesOrder(salesOrderId int) (models.SalesOrder, error) {
	db := storage.GetDB()

	var so models.SalesOrder

	query := `
	SELECT SO.id, 
	CASE SO.status
		WHEN 'NQ' THEN 'Novo orçamento'
		WHEN 'NO' THEN 'Novo pedido'
		WHEN 'QC' THEN 'Orçamento cancelado'
		WHEN 'OC' THEN 'Pedido cancelado'
		WHEN 'OD' THEN 'Pedido concluído'
		ELSE 'Não reconhecido'
	END AS StatusType,
	SO.done,
	SO.discounted_amount,
	SO.total_value,
	SO.created_at,
	SO.updated_at,
	SO.delivery_at,
	C.name,
	C.email,
	C.phone,
	C.cep,
	C.street,
	C.neighborhood,
	C.complement,
	C.state,
	C.city,
	C.pf_pj,
	C.cpf,
	C.cnpj

	FROM sales_orders AS SO 
	INNER JOIN business_histories AS BHI ON SO.id_business_history = BHI.id
	INNER JOIN customers AS C ON BHI.id_customer = C.id

	WHERE SO.id = $1
	`

	err := db.QueryRow(query, salesOrderId).Scan(
		&so.Id,
		&so.StatusType,
		&so.Done,
		&so.DiscountedAmount,
		&so.TotalValue,
		&so.CreatedAt,
		&so.UpdatedAt,
		&so.DeliveryAt,
		&so.CustomerName,
		&so.CustomerEmail,
		&so.CustomerPhone,
		&so.Cep,
		&so.Street,
		&so.Neighborhood,
		&so.Complement,
		&so.State,
		&so.City,
		&so.PfPj,
		&so.Cpf,
		&so.Cnpj,
	)
	if err != nil {
		return so, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)")
	}

	return so, nil
}
