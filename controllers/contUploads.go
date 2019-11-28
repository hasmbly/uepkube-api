package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	 // _"github.com/jinzhu/gorm/dialects/mysql"
	 // "uepkube-api/models"
	 // "uepkube-api/db"
	//  "regexp"
	//  "uepkube-api/helpers"
	 "log"	 
)

// @Summary UploadsImages -> uep | kube | pendamping | verifikator
// @Tags UploadImages-Controller
// @Accept  mpfd
// @Produce  mpfd
// @Param key path string true "Key (string) -> uep | kube | pendamping | verifikator"
// @Param id query int true "id (int)"
// @Param photo formData file true "Uploads Files"
// @Param is_display query int false "int (int) -> 0 | 1"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /upload/images/{key} [post]
func UploadImages(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("Uploads File to : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {
		case "uep":
			return UploadUepFiles(c)
		case "kube":
			return UploadKubeFiles(c)
		case "pendamping":
			return UploadUepFiles(c)
		case "verifikator":
			return UploadUepFiles(c)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")
	}

	return nil
	
}