package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	 "strconv"
	"uepkube-api/helpers"
)

// @Summary GetUepById
// @Tags Uep-Controller
// @Accept  json
// @Produce  json
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /uep [get]
func GetUep(c echo.Context) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	
	
	var val string
	User 	:= models.Tbl_user{}
	R 		:= models.CustU{}


	/*check if query key -> "val"*/
	qk := c.QueryParams()
	for k,v := range qk {
		if k == "val" {
			return err
		} else if k == "id" {
			val = v[0]
			id,_ := strconv.Atoi(val)
			/*find kube by Nama_kube:*/
			if err := con.Where(&models.Tbl_user{Id_user:id}).First(&User).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
			/*find uep by id + join pendaming-user*/
			id = User.Id_user
			if err := con.Table("tbl_uep").Select("tbl_uep.bantuan_modal, tbl_uep.status, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_uep.id_pendamping").Where(&models.Tbl_uep{Id_uep:id}).Scan(&R).Error; gorm.IsRecordNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "Uep Not Found")
			}			
		}
	}

	r := &models.Jn{Msg: &models.UepKube{Uep: models.U{User, R}}}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

// @Summary GetPaginateUep
// @Tags Uep-Controller
// @Accept  json
// @Produce  json
// @Param uep body models.PosPagin true "Show Uep"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /uep [post]
func GetPaginateUep(c echo.Context) (err error) {	
	helpers.PaginateUep(c, &r)
	return c.JSON(http.StatusOK, r)
}

// @Summary AddUep
// @Tags Uep-Controller
// @Accept  json
// @Produce  json
// @Param uep body models.PosPagin true "Add Uep"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /uep/add [post]
func AddUep(c echo.Context) (err error) {
	helpers.PaginateUep(c, &r)
	return c.JSON(http.StatusOK, r)
}