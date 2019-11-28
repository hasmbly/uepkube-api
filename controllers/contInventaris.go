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

/*@Summary GetInventarisById
@Tags Inventaris-Controller
@Accept  json
@Produce  json
@Param id query int true "int"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /inventaris [get]*/
func GetInventaris(c echo.Context) error {
	/*prepare DB*/
	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)	

	// var val string
	// Inventaris 	:= models.Tbl_inventaris{}

	// /*check if query key -> "val"*/
	// qk := c.QueryParams()
	// for k,v := range qk {
	// 	if k == "val" {
	// 		val = v[0]
	// 		/*find inventaris by Nama_inventaris:*/
	// 		if err := con.Where("nama_inventaris LIKE ?", "%" + val + "%").First(&Inventaris).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Inventaris Not Found")
	// 		}		
	// 	} else if k == "id" {
	// 		val = v[0]
	// 		id,_ := strconv.Atoi(val)
	// 		/*find inventaris by Nama_inventaris:*/
	// 		if err := con.Where(&models.Tbl_inventaris{Id_inventaris:id}).First(&Inventaris).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Inventaris Not Found")
	// 		}			
	// 	}
	// }

	// helpers.SetMemberNameInventaris(&Kt, Inventaris)

	// r := &models.Jn{Msg: Kt}

	// defer con.Close()
	// return c.JSON(http.StatusOK, r)
	return nil
}

/*@Summary GetPaginateInventaris
@Tags Inventaris-Controller
@Accept  json
@Produce  json
@Param inventaris body models.PosPagin true "Show Inventaris Paginate"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /inventaris [post]*/
func GetPaginateInventaris(c echo.Context) (err error) {	
	// if err := helpers.PaginateInventaris(c, &r); err != nil {
	// 	return echo.ErrInternalServerError
	// }	
	// return c.JSON(http.StatusOK, r)
	return nil
}

/*@Summary AddInventaris
@Tags Inventaris-Controller
@Accept  json
@Produce  json
@Param inventaris body models.Tbl_kube true "Add Inventaris"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /inventaris/add [post]*/
func AddInventaris(c echo.Context) (err error) {
	// inventaris := &models.Tbl_inventaris{}

	// if err := c.Bind(inventaris); err != nil {
	// 	return err
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Create(&inventaris).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Store Data"}
	// return c.JSON(http.StatusOK, r)
	return nil
}

/*@Summary UpdateInventaris
@Tags Inventaris-Controller
@Accept  json
@Produce  json
@Param inventaris body models.Tbl_kube true "Update Inventaris"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /inventaris [put]*/
func UpdateInventaris(c echo.Context) (err error) {
	// inventaris := &models.Tbl_inventaris{}

	// if err := c.Bind(inventaris); err != nil {
	// 	return err
	// }

	// if inventaris.Id_inventaris == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Model(&models.Tbl_inventaris{}).UpdateColumns(&inventaris).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

	// if err := con.Table("tbl_inventaris").Where("id_inventaris = ?",inventaris.Id_inventaris).UpdateColumn("status", inventaris.Status).Error; err != nil {return echo.ErrInternalServerError}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Update Data"}
	// return c.JSON(http.StatusOK, r)
	return nil	
}

/*@Summary DeleteInventaris
@Tags Inventaris-Controller
@Accept  json
@Produce  json
@Param id path int true "Delete Inventaris by id"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /inventaris/{id} [post]*/
func DeleteInventaris(c echo.Context) (err error) {
	// id, _ := strconv.Atoi(c.Param("id"))

	// if id == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	// }

	// inventaris := &models.Tbl_inventaris{}
	// inventaris.Id_inventaris = id

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Delete(&inventaris).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Delete Data"	}
	// return c.JSON(http.StatusOK, r)	
	return nil
}