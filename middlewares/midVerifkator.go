package middlewares

import (
	"net/http"
	"github.com/labstack/echo"
	"uepkube-api/helpers"
)

var VerifikatorRoles = []string{"VERIFIKATOR", "ADMINISTRATOR"}

// Middleware Check UEP roles
func CheckVerifikatoRoles(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helpers.GetLoggedUser(c,"roles")
		ch := IsPresentRolesVer(role)
		if ch == false {
			return echo.NewHTTPError(http.StatusUnauthorized, "Sorry, You can't access this resource")
		}		
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func IsPresentRolesVer(s string) (r bool){
	for _,n := range VerifikatorRoles {
		if s == n {
			return true
		}
	}
	return false
}