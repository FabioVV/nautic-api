package utils

import (
	"net/http"
	"strings"
)

/*
 Definition of user errors
 Used to make error returning easier in the API
*/

type UserError struct {
	HttpErrCode int
	ErrCode     string
	Message     string
}

var UserErrsMap = map[string]UserError{
	"email_unique": {
		HttpErrCode: http.StatusBadRequest,
		Message:     "Email already registered",
		ErrCode:     "u1",
	},
	"crypto/bcrypt: hashedPassword is not the hash of the given password": {
		HttpErrCode: http.StatusUnauthorized,
		Message:     "Invalid credentials",
		ErrCode:     "u2",
	},
}

func CheckForUserError(errToCheckFor string, err error) (UserError, bool) {
	if strings.Contains(err.Error(), errToCheckFor) {
		return UserErrsMap[errToCheckFor], true
	}
	return UserError{}, false
}

func GetUserError(errToCheckFor string) (UserError, bool) {
	if v, ok := UserErrsMap[errToCheckFor]; ok {
		return v, true
	}
	return UserError{}, false
}
