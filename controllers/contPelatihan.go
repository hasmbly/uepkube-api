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

// @Summary GetPelatihanById
// @Tags Pelatihan-Controller
// @Accept  json
// @Produce  json
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /pelatihan [get]
func GetPelatihan(c echo.Context) error {
	/*prepare DB*/
	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)	

	// var val string
	// Pelatihan 	:= models.Tbl_pelatihan{}

	// /*check if query key -> "val"*/
	// qk := c.QueryParams()
	// for k,v := range qk {
	// 	if k == "val" {
	// 		val = v[0]
	// 		/*find pelatihan by Nama_pelatihan:*/
	// 		if err := con.Where("nama_pelatihan LIKE ?", "%" + val + "%").First(&Pelatihan).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Pelatihan Not Found")
	// 		}		
	// 	} else if k == "id" {
	// 		val = v[0]
	// 		id,_ := strconv.Atoi(val)
	// 		/*find pelatihan by Nama_pelatihan:*/
	// 		if err := con.Where(&models.Tbl_pelatihan{Id_pelatihan:id}).First(&Pelatihan).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Pelatihan Not Found")
	// 		}			
	// 	}
	// }

	// helpers.SetMemberNamePelatihan(&Kt, Pelatihan)

	// r := &models.Jn{Msg: Kt}

	// defer con.Close()
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary GetPaginatePelatihan
// @Tags Pelatihan-Controller
// @Accept  json
// @Produce  json
// @Param pelatihan body models.PosPagin true "Show Pelatihan Paginate"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /pelatihan [post]
func GetPaginatePelatihan(c echo.Context) (err error) {	
	// if err := helpers.PaginatePelatihan(c, &r); err != nil {
	// 	return echo.ErrInternalServerError
	// }	
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary AddPelatihan
// @Tags Pelatihan-Controller
// @Accept  json
// @Produce  json
// @Param pelatihan body models.Tbl_kube true "Add Pelatihan"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /pelatihan/add [post]
func AddPelatihan(c echo.Context) (err error) {
	// pelatihan := &models.Tbl_pelatihan{}

	// if err := c.Bind(pelatihan); err != nil {
	// 	return err
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Create(&pelatihan).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Store Data"}
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary UpdatePelatihan
// @Tags Pelatihan-Controller
// @Accept  json
// @Produce  json
// @Param pelatihan body models.Tbl_kube true "Update Pelatihan"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /pelatihan [put]
func UpdatePelatihan(c echo.Context) (err error) {
	// pelatihan := &models.Tbl_pelatihan{}

	// if err := c.Bind(pelatihan); err != nil {
	// 	return err
	// }

	// if pelatihan.Id_pelatihan == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Model(&models.Tbl_pelatihan{}).UpdateColumns(&pelatihan).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

	// if err := con.Table("tbl_pelatihan").Where("id_pelatihan = ?",pelatihan.Id_pelatihan).UpdateColumn("status", pelatihan.Status).Error; err != nil {return echo.ErrInternalServerError}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Update Data"}
	// return c.JSON(http.StatusOK, r)
	return nil	
}

// @Summary DeletePelatihan
// @Tags Pelatihan-Controller
// @Accept  json
// @Produce  json
// @Param id path int true "Delete Pelatihan by id"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /pelatihan/{id} [post]
func DeletePelatihan(c echo.Context) (err error) {
	// id, _ := strconv.Atoi(c.Param("id"))

	// if id == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	// }

	// pelatihan := &models.Tbl_pelatihan{}
	// pelatihan.Id_pelatihan = id

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Delete(&pelatihan).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Delete Data"	}
	// return c.JSON(http.StatusOK, r)	
	return nil
}