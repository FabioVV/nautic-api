package utils

import "net/http"

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
}

func CheckForUserError(errToCheckFor string, err error) (UserError, bool) {
	if errMsg, ok := UserErrsMap[errToCheckFor]; ok {
		return errMsg, true
	}
	return UserError{}, false
}
