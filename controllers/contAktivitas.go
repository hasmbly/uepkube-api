package controllers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	//  _"github.com/jinzhu/gorm/dialects/mysql"
	 // "uepkube-api/models"
	//  "uepkube-api/db"
	//  "strconv"
	//  "uepkube-api/helpers"
)

// @Summary GetAktivitasById
// @Tags Aktivitas-Controller
// @Accept  json
// @Produce  json
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /aktivitas [get]
func GetAktivitas(c echo.Context) error {
	/*prepare DB*/
	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)	

	// var val string
	// Aktivitas 	:= models.Tbl_aktivitas{}

	// /*check if query key -> "val"*/
	// qk := c.QueryParams()
	// for k,v := range qk {
	// 	if k == "val" {
	// 		val = v[0]
	// 		/*find aktivitas by Nama_aktivitas:*/
	// 		if err := con.Where("nama_aktivitas LIKE ?", "%" + val + "%").First(&Aktivitas).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Aktivitas Not Found")
	// 		}		
	// 	} else if k == "id" {
	// 		val = v[0]
	// 		id,_ := strconv.Atoi(val)
	// 		/*find aktivitas by Nama_aktivitas:*/
	// 		if err := con.Where(&models.Tbl_aktivitas{Id_aktivitas:id}).First(&Aktivitas).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Aktivitas Not Found")
	// 		}			
	// 	}
	// }

	// helpers.SetMemberNameAktivitas(&Kt, Aktivitas)

	// r := &models.Jn{Msg: Kt}

	// defer con.Close()
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary GetPaginateAktivitas
// @Tags Aktivitas-Controller
// @Accept  json
// @Produce  json
// @Param aktivitas body models.PosPagin true "Show Aktivitas Paginate"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /aktivitas [post]
func GetPaginateAktivitas(c echo.Context) (err error) {	
	// if err := helpers.PaginateAktivitas(c, &r); err != nil {
	// 	return echo.ErrInternalServerError
	// }	
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary AddAktivitas
// @Tags Aktivitas-Controller
// @Accept  json
// @Produce  json
// @Param aktivitas body models.Tbl_kube true "Add Aktivitas"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /aktivitas/add [post]
func AddAktivitas(c echo.Context) (err error) {
	// aktivitas := &models.Tbl_aktivitas{}

	// if err := c.Bind(aktivitas); err != nil {
	// 	return err
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Create(&aktivitas).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Store Data"}
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary UpdateAktivitas
// @Tags Aktivitas-Controller
// @Accept  json
// @Produce  json
// @Param aktivitas body models.Tbl_kube true "Update Aktivitas"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /aktivitas [put]
func UpdateAktivitas(c echo.Context) (err error) {
	// aktivitas := &models.Tbl_aktivitas{}

	// if err := c.Bind(aktivitas); err != nil {
	// 	return err
	// }

	// if aktivitas.Id_aktivitas == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Model(&models.Tbl_aktivitas{}).UpdateColumns(&aktivitas).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

	// if err := con.Table("tbl_aktivitas").Where("id_aktivitas = ?",aktivitas.Id_aktivitas).UpdateColumn("status", aktivitas.Status).Error; err != nil {return echo.ErrInternalServerError}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Update Data"}
	// return c.JSON(http.StatusOK, r)
	return nil	
}

// @Summary DeleteAktivitas
// @Tags Aktivitas-Controller
// @Accept  json
// @Produce  json
// @Param id path int true "Delete Aktivitas by id"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /aktivitas/{id} [post]
func DeleteAktivitas(c echo.Context) (err error) {
	// id, _ := strconv.Atoi(c.Param("id"))

	// if id == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	// }

	// aktivitas := &models.Tbl_aktivitas{}
	// aktivitas.Id_aktivitas = id

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Delete(&aktivitas).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Delete Data"	}
	// return c.JSON(http.StatusOK, r)	
	return nil
}