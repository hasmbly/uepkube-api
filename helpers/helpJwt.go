package helpers

import (
	"github.com/labstack/echo"
	 "log"
	"github.com/dgrijalva/jwt-go"
	"uepkube-api/models"
)

func GetLoggedUser(c echo.Context, s string) string {
	var detail string
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*models.Claims)
	switch s {
	case "nama":
		detail = claims.Name
	case "roles":
		detail = claims.Roles			
	}
	log.Println("User's Roles : ",detail)
	return detail
}