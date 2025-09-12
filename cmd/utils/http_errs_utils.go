package utils

import (
	"net/http"
	"strings"
)

/*
 Definition of user errors
 Used to make error returning easier in the API
 ! CURRENTLY NOT WIDELY USED !
*/

type UserError struct {
	HttpErrCode int
	ErrCode     string
	Message     string
	Status      string `json:"status"`
}

var ErrsMap = map[string]UserError{
	"unique_email": {
		HttpErrCode: http.StatusBadRequest,
		Message:     "Email already registered",
		ErrCode:     "u1",
		Status:      "error",
	},
	"unique_type": {
		HttpErrCode: http.StatusBadRequest,
		Message:     "Email already registered",
		ErrCode:     "u1",
		Status:      "error",
	},
	"crypto/bcrypt: hashedPassword is not the hash of the given password": {
		HttpErrCode: http.StatusUnauthorized,
		Message:     "Invalid credentials",
		ErrCode:     "u2",
		Status:      "error",
	},
}

func CheckForUserError(errToCheckFor string, err error) (UserError, bool) {
	if strings.Contains(err.Error(), errToCheckFor) {
		return ErrsMap[errToCheckFor], true
	}
	return UserError{}, false
}

func GetUserError(errToCheckFor string) (UserError, bool) {
	if v, ok := ErrsMap[errToCheckFor]; ok {
		return v, true
	}
	return UserError{}, false
}
