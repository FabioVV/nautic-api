package repositories

import (
	"database/sql"
	"fmt"
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
	LEFT JOIN sales_orders_itens AS SOI ON SO.id = SOI.id_sales_order

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
		&so.OrderBoatPrice,
		&so.OrderEnginePrice,
		&so.OrderEngineId,
		&so.OrderBoatId,
		&so.OrderEngineModel,
		&so.OrderBoatModel,
	)
	if err != nil {
		return so, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)"+err.Error())
	}

	return so, nil
}

func GetSalesOrderItens(salesOrderId int) ([]models.SalesOrderItem, error) {
	db := storage.GetDB()
	var sos []models.SalesOrderItem

	query := `
	SELECT SOS.id, SOS.id_accessory, SOS.discount, A.model, SOS.qty, SOS.unit_price, SUM(SOS.unit_price * SOS.qty) OVER () AS total_price_itens
	FROM sales_orders_itens AS SOS 
	LEFT JOIN accessories AS A ON SOS.id_accessory = A.id

	WHERE SOS.id_sales_order = $1 
	AND SOS.id_boat IS NULL 
	AND SOS.id_engine IS NULL
	`

	rows, err := db.Query(query, salesOrderId)

	if err != nil {
		if err == sql.ErrNoRows {
			return sos, echo.NewHTTPError(http.StatusNotFound, "Sales order items not found")
		}
		return sos, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve sales order items"+err.Error())
	}

	for rows.Next() {
		var curSos models.SalesOrderItem

		if err := rows.Scan(&curSos.Id, &curSos.AccessoryId, &curSos.Discount, &curSos.Model, &curSos.Quantity, &curSos.UnitPrice, &curSos.TotalPriceItens); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		sos = append(sos, curSos)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sos, nil
}
