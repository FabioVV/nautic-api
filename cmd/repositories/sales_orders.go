package repositories

import (
	"database/sql"
	"fmt"
	"io"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func GetSalesOrderQuote(salesOrderId string) (models.SalesOrder, error) {
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
	SO.uuid,
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
	(SELECT BB.model FROM sales_orders_itens AS SOIi INNER JOIN boats AS BB ON SOIi.id_boat = BB.id WHERE id_sales_order = SO.id AND SOIi.id_boat IS NOT NULL ) AS order_boat_model,
	(SELECT SUM(SOIi.unit_price * SOIi.qty) FROM sales_orders_itens AS SOIi WHERE SOIi.id_sales_order = SO.id AND SOIi.id_boat IS NULL AND SOIi.id_engine IS NULL ) AS total_price_itens

	FROM sales_orders AS SO 
	INNER JOIN business_histories AS BHI ON SO.id_business_history = BHI.id
	INNER JOIN customers AS C ON BHI.id_customer = C.id
	INNER JOIN users AS U ON BHI.id_user = U.id
	LEFT JOIN sales_orders_itens AS SOI ON SO.id = SOI.id_sales_order

	WHERE SO.uuid = $1
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
		&so.Uuid,
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
		&so.TotalItensPrice,
	)
	if err != nil {
		return so, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)")
	}

	return so, nil
}

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
	SO.uuid,
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
	(SELECT BB.model FROM sales_orders_itens AS SOIi INNER JOIN boats AS BB ON SOIi.id_boat = BB.id WHERE id_sales_order = SO.id AND SOIi.id_boat IS NOT NULL ) AS order_boat_model,
	(SELECT SUM(SOIi.unit_price * SOIi.qty) FROM sales_orders_itens AS SOIi WHERE SOIi.id_sales_order = SO.id AND SOIi.id_boat IS NULL AND SOIi.id_engine IS NULL ) AS total_price_itens

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
		&so.Uuid,
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
		&so.TotalItensPrice,
	)
	if err != nil {
		return so, echo.NewHTTPError(http.StatusInternalServerError, "Internal server error (db)")
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

func ChangeSalesOrderFileType(salesOrderId int, fileId int, _type *models.UpdateSalesOrderFileType) error {
	db := storage.GetDB()

	salesOr, err := GetSalesOrder(salesOrderId)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	query := "UPDATE sales_orders_files SET type = $1 WHERE id = $2 AND id_sales_order = $3"
	_, err = db.Exec(query, _type.Type, fileId, salesOrderId)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update sales order file type")
	}

	return nil
}

func UploadSalesOrderFile(c echo.Context, id int) error {
	db := storage.GetDB()
	salesOr, err := GetSalesOrder(int(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get sales order")
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get file from request")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open file")
	}
	defer src.Close()

	uploadDir := filepath.Join(".", "uploads", "sales_orders", fmt.Sprintf("%d", id))
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create upload directory")
	}

	ts := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS
	fname := filepath.Base(file.Filename)
	fname = strings.ReplaceAll(fname, " ", "_")
	dstName := fmt.Sprintf("so_%d_%s_%s", id, ts, fname)
	dstPath := filepath.Join(uploadDir, dstName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create destination file")
	}
	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file")
	}

	query := "INSERT INTO sales_orders_files (path, id_sales_order) VALUES ($1, $2)"

	_, err = db.Exec(query, dstPath, id)
	if err != nil {
		return err
	}

	return nil
}

func GetSalesOrderFiles(salesOrderId int) ([]models.SalesOrderFile, error) {
	db := storage.GetDB()
	var soFiles []models.SalesOrderFile
	query := `
	SELECT SOF.id, SOF.path, SOF.type
	FROM sales_orders_files AS SOF
	WHERE SOF.id_sales_order = $1 AND SOF.soft_deleted = 'N'
	ORDER BY SOF.id
	`

	rows, err := db.Query(query, salesOrderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return soFiles, echo.NewHTTPError(http.StatusNotFound, "Sales order files not found")
		}
		return soFiles, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve sales order files"+err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var curSof models.SalesOrderFile
		if err := rows.Scan(&curSof.Id, &curSof.Path, &curSof.Type); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		curSof.Path = strings.ReplaceAll(curSof.Path, "\\", "/") // for windows paths
		curSof.Path = "http://127.0.0.1:8080/" + curSof.Path
		soFiles = append(soFiles, curSof)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return soFiles, nil
}

func RemoveSalesOrderFile(salesOrderId int, fileId int) error {
	db := storage.GetDB()
	// var filePath string
	// query := "SELECT path FROM sales_orders_files WHERE id = $1 AND id_sales_order = $2"
	// err := db.QueryRow(query, fileId, salesOrderId).Scan(&filePath)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return echo.NewHTTPError(http.StatusNotFound, "Sales order file not found")
	// 	}
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve sales order file"+err.Error())
	// }

	// err = os.Remove(filePath)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Could not remove sales order file"+err.Error())
	// }

	salesOr, err := GetSalesOrder(salesOrderId)
	if err != nil {
		return err
	}

	if *salesOr.StatusType == "Orçamento cancelado" || *salesOr.StatusType == "Pedido cancelado" || *salesOr.StatusType == "Pedido concluído" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Sales order cant be updated")
	}

	query := "UPDATE sales_orders_files SET soft_deleted = 'Y' WHERE id = $1 AND id_sales_order = $2"
	_, err = db.Exec(query, fileId, salesOrderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete sales order file"+err.Error())
	}

	return nil
}
