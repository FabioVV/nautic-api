package testutils

import ("nautic/cmd/storage"  "time"
"math/rand")

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStrSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generate_users(how_many int) {
	storage.InitDB()
	defer storage.CloseDB()

	db := storage.GetDB()

	var users []models.User
	query := `INSERT INTO users(
		name, email, phone, password_hash)
		VALUES ($1, $2, $3, $4)
	`
	rows, err := db.Query(query, pagenumber, paginationRange)

	if err != nil {
		if err == sql.ErrNoRows {
			return users, 0, echo.NewHTTPError(http.StatusNotFound, "Users not found")
		}
		return users, 0, echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve users")
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
