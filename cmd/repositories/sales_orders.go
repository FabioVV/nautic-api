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
	SELECT DISTINCT SO.id, 
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
	U.name,
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
	C.cnpj,
	(SELECT unit_price FROM sales_orders_itens WHERE id_sales_order = SO.id AND id_boat IS NOT NULL ) AS boat_price,
	(SELECT unit_price FROM sales_orders_itens WHERE id_sales_order = SO.id AND id_engine IS NOT NULL ) AS engine_price,
	(SELECT id_engine FROM sales_orders_itens WHERE id_sales_order = SO.id AND id_engine IS NOT NULL ) AS order_engine_id,
	(SELECT id_boat FROM sales_orders_itens WHERE id_sales_order = SO.id AND id_boat IS NOT NULL ) AS order_boat_id,
	(SELECT EE.model FROM sales_orders_itens AS SOIi INNER JOIN engines AS EE ON SOIi.id_engine = EE.id WHERE id_sales_order = SO.id AND SOIi.id_engine IS NOT NULL ) AS order_engine_model,
	(SELECT BB.model FROM sales_orders_itens AS SOIi INNER JOIN boats AS BB ON SOIi.id_boat = BB.id WHERE id_sales_order = SO.id AND SOIi.id_boat IS NOT NULL ) AS order_boat_model

	FROM sales_orders AS SO 
	INNER JOIN business_histories AS BHI ON SO.id_business_history = BHI.id
	INNER JOIN customers AS C ON BHI.id_customer = C.id
	INNER JOIN users AS U ON BHI.id_user = U.id
	INNER JOIN sales_orders_itens AS SOI ON SO.id = SOI.id_sales_order

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
		&so.SellerName,
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
		&so.OrderBoatId,
		&so.OrderBoatModel,
		&so.OrderEngineId,
		&so.OrderEngineModel,
		&so.OrderBoatPrice,
		&so.OrderEnginePrice,
	)
	if err != nil {
		return so, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)"+err.Error())
	}

	return so, nil
}
