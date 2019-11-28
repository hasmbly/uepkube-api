package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	//  _"github.com/jinzhu/gorm/dialects/mysql"
	//  "uepkube-api/models"
	//  "uepkube-api/db"
	//  "regexp"
	//  "uepkube-api/helpers"
	 "log"
)

// @Summary GetDetails -> uep | kube | pendamping | verifikator
// @Tags Users-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator"
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /users/{key} [get]
func GetUsers(c echo.Context) error {
	key 	:= c.Param("key")

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return GetUep(c)
	case "kube":
		return GetKube(c)
	case "pendamping":
		return GetPendamping(c)
	case "verifikator":
		return GetVerifikator(c)
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")
	}

	return nil
}

// @Summary Add -> uep | kube | pendamping | verifikator
// @Tags Users-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator "
// @Param uep body models.Dummy true "Add UsersUepKubePendamping"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /users/add/{key} [post]
func AddUsers(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return AddUep(c)
	case "kube":
		return AddKube(c)
	case "pendamping":
		return AddPendamping(c)
	case "verifikator":
		return AddVerifikator(c)				
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}

// @Summary Update -> uep | kube | pendamping | verifikator
// @Tags Users-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator "
// @Param uep body models.Dummy true "Update UsersUepKubePendamping"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /users/{key} [put]
func UpdateUsers(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return UpdateUep(c)
	case "kube":
		return UpdateKube(c)
	case "pendamping":
		return UpdatePendamping(c)
	case "verifikator":
		return UpdateVerifikator(c)				
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}

// @Summary Delete -> uep | kube | pendamping | verifikator
// @Tags Users-Controller
// @Accept  json
// @Produce  json
// @Param key path string true "Key -> uep | kube | pendamping | verifikator "
// @Param id path int true "Id -> uep | kube | pendamping | verifikator "
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /users/{key}/{id}  [post]
func DeleteUsers(c echo.Context) (err error) {
	key 	:= c.Param("key")

	log.Println("key : ", key)

	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill key")
	}

	switch key {

	case "uep":
		return DeleteUep(c)
	case "kube":
		return DeleteKube(c)
	case "pendamping":
		// return DeletePendamping(c)
	case "verifikator":
		// return DeleteVerifikator(c)				
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "please, choose the right key")

	}

	return nil
}