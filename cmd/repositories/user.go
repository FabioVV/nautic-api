package repositories

import (
	"database/sql"
	"fmt"
	"nautic/cmd/audit"
	"nautic/cmd/storage"
	"nautic/cmd/utils"
	"nautic/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetUserRoles(id int, id_user int) ([]string, error) {
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

	_ = audit.InsertLog(id_user, "/signin", "GET", fmt.Sprintf("GET request on user roles for login"), query)

	return roles, nil
}

func GetUserPermissions__(id int, id_user int) ([]models.UserPermission, error) {
	db := storage.GetDB()

	query := `

	SELECT
		P.id,
		UP.id,
		P.module,
		P.description,
		CASE WHEN UP.id_permission IS NULL THEN FALSE ELSE TRUE END AS has_permission
	FROM permissions AS P

	LEFT JOIN user_permissions AS UP ON UP.id_permission = P.id AND UP.id_user = $1

	ORDER BY P.module
	`

	var perms []models.UserPermission

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return perms, nil
		}
		return perms, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user permissions")
	}

	for rows.Next() {
		var curPerm models.UserPermission
		rows.Scan(&curPerm.PermissionId, &curPerm.UserPermissionID, &curPerm.Module, &curPerm.Description, &curPerm.HasPermission)
		perms = append(perms, curPerm)
	}

	if rows.Err() != nil {
		return perms, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user permissions")
	}

	_ = audit.InsertLog(id_user, "/signin", "GET", fmt.Sprintf("GET request on user roles permissions for login"), query)

	return perms, nil
}

func RemoveUserPermission(idUser int, idPerm int) error {
	db := storage.GetDB()

	query := `DELETE FROM user_permissions WHERE id_user = $1 AND id_permission = $2`

	_, err := db.Exec(query, idUser, idPerm)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(idUser, "/users/:id/permissions/:id_perm", "DELETE", fmt.Sprintf("DELETE request on user of id %d to remove permission", idUser), query)

	return nil
}

func InsertUserPermission(idUser int, idPerm int) error {
	db := storage.GetDB()

	query := `INSERT INTO user_permissions(id_user, id_permission) VALUES ($1, $2)`

	_, err := db.Exec(query, idUser, idPerm)
	if err != nil {
		return err
	}
	_ = audit.InsertLog(idUser, "/users/:id/permissions/:id_perm", "POST", fmt.Sprintf("POST request on user of id %d roles permissions for login", idUser), query)

	return nil
}

func GetUserPermissions(id int, id_user int) ([]string, error) {
	db := storage.GetDB()

	query := `SELECT P.code
	FROM user_permissions AS UP
	INNER JOIN users AS U ON UP.id_user = U.id AND U.active = 'Y'
	LEFT JOIN permissions AS P ON UP.id_permission = P.id

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

	_ = audit.InsertLog(id_user, "/users/:id/permissions", "GET", fmt.Sprintf("GET request on user of id %d roles permissions for login", id), query)

	return permissions, nil
}

func GetUser(id int, id_user int) (models.User, error) {
	db := storage.GetDB()

	var user models.User
	query := `SELECT id, name, email, active, phone, created_at, updated_at FROM users WHERE id = $1 AND active = 'Y'`

	if err := db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.Email, &user.Active, &user.Phone, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return user, echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return user, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user")
	}

	_ = audit.InsertLog(id_user, "/users/:id", "GET", fmt.Sprintf("GET request on user of id %d ", id), query)

	return user, nil
}

func GetUsers(id_user int, pagenum string, limitPerPage string, name string, email string, active string) ([]models.User, int, error) {
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
	ORDER BY U.id, U.name
	LIMIT $%d OFFSET $%d
	`, where, limitArgPos, offsetArgPos)

	rows, err := db.Query(query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return users, 0, echo.NewHTTPError(http.StatusNotFound, "Users not found")
		}
		return users, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users"+err.Error())
	}

	queryTotalRecords := fmt.Sprintf(`
	SELECT COUNT(1)
	FROM users AS U
	%s
	`, where)
	//println(queryTotalRecords)

	rowsCount := db.QueryRow(queryTotalRecords, args[:len(args)-2]...) // slice to remove the limit and offset args, they are not needed here
	numRecords := 0
	rowsCount.Scan(&numRecords)

	for rows.Next() {
		var curUser models.User
		rows.Scan(&curUser.Id, &curUser.Name, &curUser.Email, &curUser.Active, &curUser.CreatedAt, &curUser.UpdatedAt)
		users = append(users, curUser)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	_ = audit.InsertLog(id_user, "/users", "GET", fmt.Sprintf("GET request on users"), query)

	return users, numRecords, nil
}

func InsertUser(user *models.CreateUserRequest, id_user int) error {
	db := storage.GetDB()

	if errMsg, ok := utils.IsGoodPassword(user.Password); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (name, email, phone, password_hash) VALUES ($1, $2, $3, $4) RETURNING ID"

	var id int64 = 0
	err = db.QueryRow(query, user.Name, user.Email, user.Phone, hashedPassword).Scan(&id)
	if err != nil {
		if errU, ok := utils.CheckForError("unique_email", err); ok {
			return echo.NewHTTPError(errU.HttpErrCode, errU)
		}
		return err
	}

	rolePermissions, err := GetRolePermissionsByName(user.Role, id_user)
	if err != nil {
		return err
	}

	if len(rolePermissions) > 0 && id != 0 || user.Role == "Admin" {
		for _, k := range rolePermissions {
			query := "INSERT INTO user_permissions (id_user, id_permission) VALUES ($1, $2)"
			_, err = db.Exec(query, id, k.PermissionId)
			if err != nil {
				return err
			}
		}

		role, err := GetRoleByName(user.Role, id_user)
		if err != nil {
			return err
		}

		queryRole := "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)"
		role, _ = GetRoleByName(user.Role, id_user)
		_, err = db.Exec(queryRole, id, role.Id)
		if err != nil {
			return err
		}
	}

	_ = audit.InsertLog(id_user, "/users", "POST", fmt.Sprintf("POST request on users"), query)

	return nil
}

func UpdateUser(id int, id_user int, user *models.UpdateUserRequest) error {
	db := storage.GetDB()

	_, err := GetUser(id, id_user)
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

	_ = audit.InsertLog(id_user, "/users", "UPDATE", fmt.Sprintf("UPDATE request on user of id %d", id), query)

	return nil
}

func DeactivateUser(id int, id_user int) error {
	db := storage.GetDB()

	_, err := GetUser(id, id_user)
	if err != nil {
		return err
	}

	query := `UPDATE users SET active = 'N' WHERE id = $1`

	_, err = db.Exec(query, id)
	if err != nil {
		return err
	}

	_ = audit.InsertLog(id_user, "/users", "DELETE", fmt.Sprintf("DELETE request on user of id %d", id), query)

	return nil
}
