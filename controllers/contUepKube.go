package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	 "regexp"
	 "uepkube-api/helpers"
)

// @Summary Get Uep (byNik) or Get Kube (byName)
// @Tags Lookup-Controller
// @Accept  json
// @Produce  json
// @Param val query string true "string"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /lookup/uepkube [get]
func GetUepOrKube(c echo.Context) error {
	val 	:= c.QueryParam("val")
	User 	:= models.Tbl_user{}
	Kube 	:= models.Tbl_kube{}
	R 		:= models.CustU{}
	Kt 		:= models.Ktype{}

	re := regexp.MustCompile("[0-9]+")
	errr := (re.FindAllString(val, -1))

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	/*query user*/
	if errr != nil {
		if err := con.Where(&models.Tbl_user{Nik:val}).First(&User).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
		/*find uep by id + join pendaming-user*/
		id := User.Id_user
		if err := con.Table("tbl_uep").Select("tbl_uep.bantuan_modal, tbl_uep.status, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_uep.id_pendamping").Where(&models.Tbl_uep{Id_uep:id}).Scan(&R).Error; gorm.IsRecordNotFoundError(err) {
			return echo.NewHTTPError(http.StatusNotFound, "Uep Not Found")
		}
		/*find kube by Ketua:id*/
		if err := con.Where(&models.Tbl_kube{Ketua:id}).First(&Kube).Error; gorm.IsRecordNotFoundError(err) {
			r := &models.Jn{Msg: &models.UepKube{Uep: models.U{User, R}}}
			defer con.Close()
			return c.JSON(http.StatusOK, r)
		}			
	} else if errr == nil {
		return GetKube(c)
	}

	helpers.SetMemberNameKube(&Kt, Kube)

	r := &models.Jn{Msg: &models.UepKube{Uep: models.U{User, R},Kube: Kt}}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}
