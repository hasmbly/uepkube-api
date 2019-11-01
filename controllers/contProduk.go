package controllers

import (
	// "net/http"
	"github.com/labstack/echo"
	// "github.com/jinzhu/gorm"
	//  _"github.com/jinzhu/gorm/dialects/mysql"
	 // "uepkube-api/models"
	//  "uepproduk-api/db"
	//  "strconv"
	//  "uepproduk-api/helpers"
)

// @Summary GetProdukById
// @Tags Produk-Controller
// @Accept  json
// @Produce  json
// @Param id query int true "int"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /produk [get]
func GetProduk(c echo.Context) error {
	/*prepare DB*/
	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)	

	// var val string
	// Produk 	:= models.Tbl_produk{}

	// /*check if query key -> "val"*/
	// qk := c.QueryParams()
	// for k,v := range qk {
	// 	if k == "val" {
	// 		val = v[0]
	// 		/*find produk by Nama_produk:*/
	// 		if err := con.Where("nama_produk LIKE ?", "%" + val + "%").First(&Produk).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Produk Not Found")
	// 		}		
	// 	} else if k == "id" {
	// 		val = v[0]
	// 		id,_ := strconv.Atoi(val)
	// 		/*find produk by Nama_produk:*/
	// 		if err := con.Where(&models.Tbl_produk{Id_produk:id}).First(&Produk).Error; gorm.IsRecordNotFoundError(err)  {
	// 			return echo.NewHTTPError(http.StatusNotFound, "Produk Not Found")
	// 		}			
	// 	}
	// }

	// helpers.SetMemberNameProduk(&Kt, Produk)

	// r := &models.Jn{Msg: Kt}

	// defer con.Close()
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary GetPaginateProduk
// @Tags Produk-Controller
// @Accept  json
// @Produce  json
// @Param produk body models.PosPagin true "Show Produk Paginate"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /produk [post]
func GetPaginateProduk(c echo.Context) (err error) {	
	// if err := helpers.PaginateProduk(c, &r); err != nil {
	// 	return echo.ErrInternalServerError
	// }	
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary AddProduk
// @Tags Produk-Controller
// @Accept  json
// @Produce  json
// @Param produk body models.Tbl_kube true "Add Produk"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /produk/add [post]
func AddProduk(c echo.Context) (err error) {
	// produk := &models.Tbl_produk{}

	// if err := c.Bind(produk); err != nil {
	// 	return err
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Create(&produk).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Store Data"}
	// return c.JSON(http.StatusOK, r)
	return nil
}

// @Summary UpdateProduk
// @Tags Produk-Controller
// @Accept  json
// @Produce  json
// @Param produk body models.Tbl_kube true "Update Produk"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /produk [put]
func UpdateProduk(c echo.Context) (err error) {
	// produk := &models.Tbl_produk{}

	// if err := c.Bind(produk); err != nil {
	// 	return err
	// }

	// if produk.Id_produk == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	// }

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Model(&models.Tbl_produk{}).UpdateColumns(&produk).Error; err != nil {
	// 	return echo.ErrInternalServerError
	// }

	// if err := con.Table("tbl_produk").Where("id_produk = ?",produk.Id_produk).UpdateColumn("status", produk.Status).Error; err != nil {return echo.ErrInternalServerError}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Update Data"}
	// return c.JSON(http.StatusOK, r)
	return nil	
}

// @Summary DeleteProduk
// @Tags Produk-Controller
// @Accept  json
// @Produce  json
// @Param id path int true "Delete Produk by id"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @security ApiKeyAuth
// @Router /produk/{id} [post]
func DeleteProduk(c echo.Context) (err error) {
	// id, _ := strconv.Atoi(c.Param("id"))

	// if id == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	// }

	// produk := &models.Tbl_produk{}
	// produk.Id_produk = id

	// con, err := db.CreateCon()
	// if err != nil { return echo.ErrInternalServerError }
	// con.SingularTable(true)

	// if err := con.Delete(&produk).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// defer con.Close()

	// r := &models.Jn{Msg: "Success Delete Data"	}
	// return c.JSON(http.StatusOK, r)	
	return nil
}