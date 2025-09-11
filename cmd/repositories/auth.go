package repositories

import (
	"nautic/cmd/middleware"
)

func GetPermissions() (map[string]string, error) {
	return middleware.RoutesPermissions, nil
}
