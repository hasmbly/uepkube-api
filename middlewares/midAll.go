package middlewares

import (
	"net/http"
	"github.com/labstack/echo"
	"uepkube-api/helpers"
)

var AllRoles = []string{"PENDAMPING_UEP", "VERIFIKATOR", "PENDAMPING_KUBE"}

// Middleware Check UEP roles
func CheckAllRoles(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := helpers.GetLoggedUser(c,"roles")
		ch := IsPresentRolesAll(role)
		if ch == false {
			return echo.NewHTTPError(http.StatusUnauthorized, "Sorry, You can't access this resource")
		}		
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func IsPresentRolesAll(s string) (r bool){
	for _,n := range AllRoles {
		if s == n {
			return true
		}
	}
	return false
}