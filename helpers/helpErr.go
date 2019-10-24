package helpers

import (
	"net/http"
	"github.com/labstack/echo"
	"time"
	"uepkube-api/models"
)

//NewError example
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code 	:= http.StatusInternalServerError
	loc, _ 	:= time.LoadLocation("Asia/Jakarta")
	tm 		:= time.Now().In(loc)
	var msg string

	if she, ok := err.(*echo.HTTPError); ok {
		msg = she.Message.(string)
	}

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	er := models.HTTPError{
	        Code:    code,
	        Message:  msg,
	        Times: tm,
	}
	
	c.JSON(code, er)
	c.Logger().Error(err)
}