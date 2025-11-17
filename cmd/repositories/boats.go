package repositories

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"nautic/cmd/audit"
	"nautic/cmd/storage"
	"nautic/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func InsertBoat(id_user int, acc *models.CreateBoatRequest) error {
	db := storage.GetDB()

	query := "INSERT INTO boats (model, new_used) VALUES ($1, $2)"

	_, err := db.Exec(query, acc.Model, acc.NewUsed)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats", "POST", fmt.Sprintf("POST request on boats"), query)

	return nil
}

func GetBoatAds(id_user int, boatId int) ([]models.BoatAd, error) {
	db := storage.GetDB()
	var ads []models.BoatAd

	query := `
	SELECT BA.id, BA.id_boat, BA.id_mean_communication, BA.link

	FROM boat_ads AS BA
	WHERE BA.id_boat = $1

	ORDER BY BA.id
	`

	rows, err := db.Query(query, boatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ads, echo.NewHTTPError(http.StatusNotFound, "Boats ads not found")
		}
		return ads, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats ads")
	}

	for rows.Next() {
		var curAcc models.BoatAd
		rows.Scan(&curAcc.Id, &curAcc.BoatId, &curAcc.ComMeanId, &curAcc.Link)
		ads = append(ads, curAcc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/boats/:id/ads", "GET", fmt.Sprintf("GET request on boats ads"), query)

	return ads, nil
}

func GetBoatEngines(id_user int, boatId int) ([]models.Engine, error) {
	db := storage.GetDB()
	var engs []models.Engine

	query := `
	SELECT E.id, E.model, E.type, E.weight, E.rotation, E.power, E.cylinders, E.selling_price, E.command, E.clocks, E.tempo, E.fuel_type, E.active, E.created_at, E.updated_at, E.propulsion

	FROM boat_engines AS BE
	INNER JOIN engines AS E ON BE.id_engine = E.id
	WHERE BE.id_boat = $1
	ORDER BY BE.id
	`

	rows, err := db.Query(query, boatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return engs, echo.NewHTTPError(http.StatusNotFound, "Boats engines not found")
		}
		return engs, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats engines")
	}

	for rows.Next() {
		var curAcc models.Engine
		rows.Scan(&curAcc.Id, &curAcc.Model, &curAcc.Type, &curAcc.Weight, &curAcc.Rotation, &curAcc.Power, &curAcc.Cylinders,
			&curAcc.PriceSell, &curAcc.Command, &curAcc.Clocks, &curAcc.Tempo, &curAcc.FuelType, &curAcc.Active, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Propulsion)
		engs = append(engs, curAcc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/boats/:id/engines", "GET", fmt.Sprintf("GET request on boats engines"), query)

	return engs, nil
}

func GetBoatAccessories(id_user int, boatId int) ([]models.Accessory, error) {
	db := storage.GetDB()
	var accs []models.Accessory

	query := `
	SELECT A.id, A.model, A.details, A.price_buy, A.price_sell, A.created_at, A.updated_at, A.active, AT.type

	FROM boat_accessories AS BA
	INNER JOIN accessories AS A ON BA.id_accessory = A.id
	INNER JOIN accessory_types AS AT ON A.id_accessory_type = AT.id
	WHERE BA.id_boat = $1
	ORDER BY BA.id
	`

	rows, err := db.Query(query, boatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return accs, echo.NewHTTPError(http.StatusNotFound, "Boats accessories not found")
		}
		return accs, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats accessories")
	}

	for rows.Next() {
		var curAcc models.Accessory
		rows.Scan(&curAcc.Id, &curAcc.Model, &curAcc.Details, &curAcc.PriceBuy, &curAcc.PriceSell, &curAcc.CreatedAt, &curAcc.UpdatedAt, &curAcc.Active, &curAcc.Type)
		accs = append(accs, curAcc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/boats/:id/accessories", "GET", fmt.Sprintf("GET request on boats accessories"), query)

	return accs, nil
}

func GetBoats(id_user int, pagenum string, limitPerPage string, model string, price string, id string, active string) ([]models.Boat, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var boats []models.Boat

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if model != "" {
		conds = append(conds, fmt.Sprintf("B.model ILIKE $%d", paramCount))
		args = append(args, "%"+model+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("B.active = $%d", paramCount))
		args = append(args, active)
		paramCount++
	}

	if id != "" {
		conds = append(conds, fmt.Sprintf("B.id = $%d", paramCount))
		args = append(args, id)
		paramCount++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	//append pagination range
	args = append(args, limitPerPage, offset)
	limitArgPos := paramCount
	offsetArgPos := paramCount + 1

	query := fmt.Sprintf(`
	SELECT B.id, B.model, B.selling_price, B.cost, B.itens, B.hours, B.year, B.new_used, B.cab_open, B.capacity, B.night_capacity, B.length,
	B.beam, B.draft, B.weight, B.trim, B.fuel_tank_capacity, B.active,
	B.created_at, B.updated_at

	FROM boats AS B

	%s

	ORDER BY B.id, B.model
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return boats, 0, echo.NewHTTPError(http.StatusNotFound, "Boats not found")
		}
		return boats, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boats")
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM boats AS B
	%s
	`, where)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curBoat models.Boat
		rows.Scan(&curBoat.Id, &curBoat.Model, &curBoat.PriceSell,
			&curBoat.Cost, &curBoat.Itens, &curBoat.Hours, &curBoat.Year, &curBoat.NewUsed,
			&curBoat.CabOpen, &curBoat.Capacity, &curBoat.NightCapacity, &curBoat.Length, &curBoat.Beam,
			&curBoat.Draft, &curBoat.Weight, &curBoat.Trim, &curBoat.FuelTankCapactiy, &curBoat.Active, &curBoat.CreatedAt, &curBoat.UpdatedAt)
		boats = append(boats, curBoat)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/boats", "GET", fmt.Sprintf("GET request on boats"), query)

	return boats, numRecords, nil
}

func RemoveBoatFile(id_user int, id int, id_file int) error {
	db := storage.GetDB()
	query := `UPDATE boat_files SET soft_deleted = 'Y' WHERE id = $1 AND id_boat = $2`

	_, err := db.Exec(query, id_file, id)
	if err != nil {
		return err
	}
	_ = audit.InsertLog(id_user, "/boats/:id/files/:id_file", "DELETE", fmt.Sprintf("DELETE request on boat file of id %d", id), query)

	return nil
}

func GetBoatFiles(id_user int, id int) ([]models.BoatFile, error) {
	db := storage.GetDB()
	var files []models.BoatFile

	query := `
	SELECT BF.id, BF.path 
	FROM boat_files AS BF
	WHERE BF.id_boat = $1 AND BF.soft_deleted = 'N'
	ORDER BY BF.id
	`
	rows, err := db.Query(query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return files, echo.NewHTTPError(http.StatusNotFound, "Boat files not found")
		}
		return files, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boat files")
	}

	for rows.Next() {
		var curFile models.BoatFile
		rows.Scan(&curFile.Id, &curFile.Path)

		*curFile.Path = strings.ReplaceAll(*curFile.Path, "\\", "/") // for windows paths
		*curFile.Path = "http://127.0.0.1:8080/" + *curFile.Path
		files = append(files, curFile)
	}

	_ = audit.InsertLog(id_user, "/boats/:id/files", "GET", fmt.Sprintf("GET request on boats files"), query)

	return files, nil
}

func GetBoat(id_user int, id int) (models.Boat, error) {
	db := storage.GetDB()

	var curBoat models.Boat

	query := `
	SELECT B.id, B.model, B.selling_price, B.cost, B.itens, B.hours, B.year, B.new_used, B.cab_open, B.capacity, B.night_capacity, B.length,
	B.beam, B.draft, B.weight, B.trim, B.fuel_tank_capacity, B.active,
	B.created_at, B.updated_at

	FROM boats AS B

	WHERE B.id = $1

	ORDER BY B.id, B.model
	`

	if err := db.QueryRow(query, id).Scan(&curBoat.Id, &curBoat.Model, &curBoat.PriceSell,
		&curBoat.Cost, &curBoat.Itens, &curBoat.Hours, &curBoat.Year, &curBoat.NewUsed,
		&curBoat.CabOpen, &curBoat.Capacity, &curBoat.NightCapacity, &curBoat.Length, &curBoat.Beam,
		&curBoat.Draft, &curBoat.Weight, &curBoat.Trim, &curBoat.FuelTankCapactiy, &curBoat.Active, &curBoat.CreatedAt, &curBoat.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return curBoat, echo.NewHTTPError(http.StatusNotFound, "Negotiation not found")
		}
		return curBoat, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve boat")
	}

	_ = audit.InsertLog(id_user, "/boats/:id", "GET", fmt.Sprintf("GET request on boat"), query)

	return curBoat, nil
}

func RemoveBoatAccessory(id_user int, id int, id_acc int) error {
	db := storage.GetDB()

	_, err := GetAccessory(id_acc, id_user)
	if err != nil {
		return err
	}

	query := `DELETE FROM boat_accessories WHERE id_accessory = $1 AND id_boat = $2`

	_, err = db.Exec(query, id_acc, id)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/accessories/:id_acc", "DELETE", fmt.Sprintf("DELETE request on boat accessory of id %d", id), query)

	return nil
}

func RemoveBoatAd(id_user int, id int, id_mean int) error {
	db := storage.GetDB()

	// _, err := GetBoatAd(id_mean)
	// if err != nil {
	// 	return err
	// }

	query := `DELETE FROM boat_ads WHERE id_mean_communication = $1 AND id_boat = $2`

	_, err := db.Exec(query, id_mean, id)
	if err != nil {
		return err
	}
	_ = audit.InsertLog(id_user, "/boats/:id/ads/:id_ad", "DELETE", fmt.Sprintf("DELETE request on boat ad of id %d", id), query)

	return nil
}

func RemoveBoatEngine(id_user int, id int, id_eng int) error {
	db := storage.GetDB()

	_, err := GetEngine(id_eng)
	if err != nil {
		return err
	}

	query := `DELETE FROM boat_engines WHERE id_engine = $1 AND id_boat = $2`

	_, err = db.Exec(query, id_eng, id)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/engines/:id_eng", "DELETE", fmt.Sprintf("DELETE request on boat engine of id %d", id), query)

	return nil
}

func InsertEngineAccessory(id_user int, id int, id_eng int) error {
	db := storage.GetDB()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM boat_engines WHERE id_boat = $1 AND id_engine = $2)`
	if err := db.QueryRow(checkQuery, id, id_eng).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "engine already linked with the boat")
	}

	query := "INSERT INTO boat_engines (id_boat, id_engine) VALUES ($1, $2)"

	_, err := db.Exec(query, id, id_eng)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/engines", "POST", fmt.Sprintf("POST request on boat of id %d engine ", id), query)

	return nil
}

func InsertBoatAd(id_user int, id int, id_mean int, bAd *models.BoatAdCreate) error {
	db := storage.GetDB()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM boat_ads WHERE id_boat = $1 AND id_mean_communication = $2)`
	if err := db.QueryRow(checkQuery, id, id_mean).Scan(&exists); err != nil {
		return err
	}

	// if exists {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "ad already linked with the boat")
	// }

	query := "INSERT INTO boat_ads (id_boat, id_mean_communication, link) VALUES ($1, $2, $3)"

	_, err := db.Exec(query, id, id_mean, bAd.Link)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/ads/:id_mean", "POST", fmt.Sprintf("POST request on boat of id %d ad ", id), query)

	return nil
}

func InsertBoatAccessory(id_user int, id int, id_acc int) error {
	db := storage.GetDB()

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM boat_accessories WHERE id_boat = $1 AND id_accessory = $2)`
	if err := db.QueryRow(checkQuery, id, id_acc).Scan(&exists); err != nil {
		return err
	}

	if exists {
		return echo.NewHTTPError(http.StatusBadRequest, "accessory already linked with the boat")
	}

	query := "INSERT INTO boat_accessories (id_boat, id_accessory) VALUES ($1, $2)"

	_, err := db.Exec(query, id, id_acc)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/accessories/:id_acc", "POST", fmt.Sprintf("POST request on boat of id %d accessory ", id), query)

	return nil
}

func UploadBoatFile(id_user int, id int, file *multipart.FileHeader) error {
	db := storage.GetDB()

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to process uploaded file")
	}
	defer src.Close()

	uploadDir := filepath.Join(".", "uploads", "boats", fmt.Sprintf("%d", id))
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create upload directory")
	}

	ts := time.Now().Format("20060102_150405") // YYYYMMDD_HHMMSS
	fname := filepath.Base(file.Filename)
	fname = strings.ReplaceAll(fname, " ", "_")
	dstName := fmt.Sprintf("boat_%d_%s_%s", id, ts, fname)
	dstPath := filepath.Join(uploadDir, dstName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create destination file")
	}
	if _, err := io.Copy(dst, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file")
	}

	query := "INSERT INTO boat_files (path, id_boat) VALUES ($1, $2)"

	_, err = db.Exec(query, dstPath, id)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id/boat-file", "POST", fmt.Sprintf("POST request on boat of id %d boat-file ", id), query)

	return nil
}

func UpdateBoat(id_user int, id int, cT *models.BoatRequest) error {
	db := storage.GetDB()

	query := `UPDATE boats SET `
	params := []interface{}{}
	paramCount := 0

	if cT.Model != nil {
		paramCount++
		query += fmt.Sprintf("model = $%d, ", paramCount)
		params = append(params, *&cT.Model)
	}

	if cT.Beam != nil {
		paramCount++
		query += fmt.Sprintf("beam = $%d, ", paramCount)
		params = append(params, *&cT.Beam)
	}

	if cT.Trim != nil {
		paramCount++
		query += fmt.Sprintf("trim = $%d, ", paramCount)
		params = append(params, *&cT.Trim)
	}

	if cT.Capacity != nil {
		paramCount++
		query += fmt.Sprintf("capacity = $%d, ", paramCount)
		params = append(params, *&cT.Capacity)
	}

	if cT.NightCapacity != nil {
		paramCount++
		query += fmt.Sprintf("night_capacity = $%d, ", paramCount)
		params = append(params, *&cT.NightCapacity)
	}

	if cT.Weight != nil {
		paramCount++
		query += fmt.Sprintf("weight = $%d, ", paramCount)
		params = append(params, *&cT.Weight)
	}

	if cT.Length != nil {
		paramCount++
		query += fmt.Sprintf("length = $%d, ", paramCount)
		params = append(params, *&cT.Length)
	}

	if cT.FuelTankCapactiy != nil {
		paramCount++
		query += fmt.Sprintf("fuel_tank_capacity = $%d, ", paramCount)
		params = append(params, *&cT.FuelTankCapactiy)
	}

	if cT.Draft != nil {
		paramCount++
		query += fmt.Sprintf("draft = $%d, ", paramCount)
		params = append(params, *&cT.Draft)
	}

	if cT.Itens != nil {
		paramCount++
		query += fmt.Sprintf("itens = $%d, ", paramCount)
		params = append(params, *&cT.Itens)
	}

	if cT.PriceSell != nil {
		paramCount++
		query += fmt.Sprintf("selling_price = $%d, ", paramCount)
		params = append(params, *&cT.PriceSell)
	}

	if cT.Hours != nil {
		paramCount++
		query += fmt.Sprintf("hours = $%d, ", paramCount)
		params = append(params, *&cT.Hours)
	}

	if cT.Year != nil {
		paramCount++
		query += fmt.Sprintf("year = $%d, ", paramCount)
		params = append(params, *&cT.Year)
	}

	if cT.NewUsed != nil {
		paramCount++
		query += fmt.Sprintf("new_used = $%d, ", paramCount)
		params = append(params, *&cT.NewUsed)
	}

	if cT.CabOpen != nil {
		paramCount++
		query += fmt.Sprintf("cab_open = $%d, ", paramCount)
		params = append(params, *&cT.CabOpen)
	}

	if len(params) == 0 {
		return nil
	}

	//Remove the trailing comma and space from the query
	query = query[:len(query)-2]

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	_, err := db.Exec(query, params...)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/boats/:id", "UPDATE", fmt.Sprintf("UPDATE request on boat of id %d boat-file ", id), query)

	return nil
}
