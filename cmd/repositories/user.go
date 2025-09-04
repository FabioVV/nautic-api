package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetUserRoles(id int) ([]string, error) {
	db := storage.GetDB()

	query := `SELECT R.name
	FROM user_roles AS UR
	INNER JOIN roles AS R ON UR.role_id = R.id
	WHERE UR.user_id = $1
	`

	var roles []string

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return roles, nil
		}
		return roles, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user")
	}

	for rows.Next() {
		var role string
		rows.Scan(&role)
		roles = append(roles, role)
	}

	if rows.Err() != nil {
		return []string{}, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user roles")
	}

	return roles, nil
}

func GetUserPermissions(id int) ([]string, error) {
	db := storage.GetDB()

	query := `SELECT UP.code
	FROM user_permissions AS UP
	INNER JOIN users AS U ON UP.id_user = U.id AND U.active = 'Y'
	WHERE UP.id_user = $1
	`

	var permissions []string

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return permissions, nil
		}
		return permissions, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user")
	}

	var permission string
	for rows.Next() {
		rows.Scan(&permission)
		permissions = append(permissions, permission)
	}

	if rows.Err() != nil {
		return []string{}, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user permissions")
	}

	return permissions, nil
}

func GetUser(id int) (models.User, error) {
	db := storage.GetDB()

	var user models.User
	query := `SELECT id, name, email, active, phone, created_at, updated_at FROM users WHERE id = $1 AND active = 'Y'`

	if err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Active, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return user, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user")
	}

	return user, nil
}

func GetUsers(pagenum string, limitPerPage string, name string, email string, active string) ([]models.User, int, error) {
	db := storage.GetDB()

	pagenumber, err := strconv.Atoi(pagenum)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users (PG1)")
	}
	limit, err := strconv.Atoi(limitPerPage)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users (PG2)")
	}

	offset := (pagenumber - 1) * limit

	var users []models.User

	conds := []string{}
	args := []interface{}{}
	paramCount := 1

	if name != "" {
		conds = append(conds, fmt.Sprintf("U.name ILIKE $%d", paramCount))
		args = append(args, "%"+name+"%")
		paramCount++
	}

	if email != "" {
		conds = append(conds, fmt.Sprintf("U.email ILIKE $%d", paramCount))
		args = append(args, "%"+email+"%")
		paramCount++
	}

	if active != "" {
		conds = append(conds, fmt.Sprintf("U.active = $%d", paramCount))
		args = append(args, active)
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
	SELECT U.id, U.name, U.email, U.active, U.created_at, U.updated_at
	FROM users AS U
	%s
	ORDER BY U.name
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return users, 0, echo.NewHTTPError(http.StatusNotFound, "Users not found")
		}
		return users, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users"+err.Error())
	}

	numRecords := 0
	for rows.Next() {
		var curUser models.User
		rows.Scan(&curUser.Id, &curUser.Name, &curUser.Email, &curUser.Active, &curUser.CreatedAt, &curUser.UpdatedAt)
		users = append(users, curUser)
		numRecords++
	}

	return users, numRecords, nil
}

func InsertUser(user *models.CreateUserRequest) error {
	db := storage.GetDB()

	if errMsg, ok := utils.IsGoodPassword(user.Password); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := "INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4)"

	_, err = db.Exec(query, user.Name, user.Email, user.Phone, hashedPassword)
	if err != nil {
		if errU, ok := utils.CheckForUserError("unique_email", err); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	return nil
}

func UpdateUser(id int, user *models.UpdateUserRequest) error {
	db := storage.GetDB()

	_, err := GetUser(id)
	if err != nil {
		return err
	}

	query := `UPDATE users SET`
	params := []interface{}{}
	paramCount := 0

	if user.Name != nil {
		paramCount++
		query += fmt.Sprintf("name = $%d, ", paramCount)
		params = append(params, *user.Name)
	}
	if user.Email != nil {
		paramCount++
		query += fmt.Sprintf("email = $%d, ", paramCount)
		params = append(params, *user.Email)
	}
	if user.Phone != nil {
		paramCount++
		query += fmt.Sprintf("phone = $%d, ", paramCount)
		params = append(params, *user.Phone)
	}
	if user.Active != nil {
		paramCount++
		query += fmt.Sprintf("active = $%d, ", paramCount)
		params = append(params, *user.Active)
	}

	if len(params) == 0 {
		return nil
	}

	//Remove the trailing comma and space from the query
	query = query[:len(query)-2]

	paramCount++
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, id)

	_, err = db.Exec(query, params...)
	if err != nil {
		return err
	}

	return nil

}

func DeactivateUser(id int) error {
	db := storage.GetDB()

	_, err := GetUser(id)
	if err != nil {
		return err
	}

	query := `UPDATE users SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil

}
