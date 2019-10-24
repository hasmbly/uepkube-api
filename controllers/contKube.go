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
// @Summary GetKubeById
// @Tags Kube-Controller
// @Accept  json
// @Produce  json
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /kube [get]
func GetKube(c echo.Context) error {
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	

	var val string
	Kube 	:= models.Tbl_kube{}

	/*check if query key -> "val"*/
	qk := c.QueryParams()
	for k,v := range qk {
		if k == "val" {
			val = v[0]
			/*find kube by Nama_kube:*/
			if err := con.Where("nama_kube LIKE ?", "%" + val + "%").First(&Kube).Error; gorm.IsRecordNotFoundError(err)  {
				return echo.NewHTTPError(http.StatusNotFound, "Kube Not Found")
			}		
		} else if k == "id" {
			val = v[0]
			id,_ := strconv.Atoi(val)
			/*find kube by Nama_kube:*/
			if err := con.Where(&models.Tbl_kube{Id_kube:id}).First(&Kube).Error; gorm.IsRecordNotFoundError(err)  {
				return echo.NewHTTPError(http.StatusNotFound, "Kube Not Found")
			}			
		}
	}

	helpers.SetMemberNameKube(&Kt, Kube)

	r := &models.Jn{Msg: &models.UepKube{Kube: Kt}}

	defer con.Close()
	return c.JSON(http.StatusOK, r)
}
