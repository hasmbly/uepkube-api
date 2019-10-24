package middlewares

import (
	"net/http"
	"github.com/labstack/echo"
	"uepkube-api/helpers"
)

var UepRoles = []string{"PENDAMPING_UEP", "VERIFIKATOR"}

// Middleware Check UEP roles
func CheckUepRoles(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helpers.GetLoggedUser(c,"roles")
		ch := IsPresentRolesUep(role)
		if ch == false {
			return echo.NewHTTPError(http.StatusUnauthorized, "Sorry, You can't access this resource")
		}		
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func IsPresentRolesUep(s string) (r bool){
	for _,n := range UepRoles {
		if s == n {
			return true
		}
	}
	return false
}