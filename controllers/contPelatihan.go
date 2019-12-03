package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/db"
	 "uepkube-api/models"
	 "strconv"
	 "uepkube-api/helpers"
	 "log"
)

/*@Summary GetPelatihanById
@Tags Pelatihan-Controller
@Accept  json
@Produce  json
@Param id query int true "int"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@Router /pelatihan [get]*/
func GetPelatihan(c echo.Context) error {
	id 		:= c.QueryParam("id")

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Pelatihan := models.Pelatihan{}
	q := con
	q = q.Table("tbl_pelatihan t1")
	q = q.Where("t1.id_pelatihan = ?", id)
	if ErrNo := q.Scan(&Pelatihan); ErrNo.Error != nil { 
		log.Println("Erro : ", ErrNo.Error)
		return echo.ErrNotFound
	}

    // get photo pelatihan
    var photo []models.Tbl_pelatihan_photo
	if err := con.Table("tbl_pelatihan_photo").Where(&models.Tbl_pelatihan_photo{Id_pelatihan: Pelatihan.Id_pelatihan}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			if photo[i].Photo != "" {
				ImageBlob := photo[i].Photo
				photo[i].Photo = "data:image/png;base64," + ImageBlob	
			}

		}
	Pelatihan.Photo = photo

	r := &models.Jn{Msg: Pelatihan}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

/*@Summary GetPaginatePelatihan
@Tags Pelatihan-Controller
@Accept  json
@Produce  json
@Param pelatihan body models.PosPagin true "Show Pelatihan Paginate"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@Router /pelatihan [post]*/
func GetPaginatePelatihan(c echo.Context) (err error) {	
	if err := helpers.PaginatePelatihan(c, &r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

/*@Summary AddPelatihan
@Tags Pelatihan-Controller
@Accept  json
@Produce  json
@Param pelatihan body models.Tbl_kube true "Add Pelatihan"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /pelatihan/add [post]*/
func AddPelatihan(c echo.Context) (err error) {
	pelatihan := &models.Pelatihan{}

	if err := c.Bind(pelatihan); err != nil {
		return err
	}

	Pelatihan := &models.Tbl_pelatihan{}
	Pelatihan = pelatihan.Tbl_pelatihan

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Create(&Pelatihan).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}
/*
@Summary UpdatePelatihan
@Tags Pelatihan-Controller
@Accept  json
@Produce  json
@Param pelatihan body models.Tbl_kube true "Update Pelatihan"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /pelatihan [put]*/
func UpdatePelatihan(c echo.Context) (err error) {
	pelatihan := &models.Pelatihan{}

	if err := c.Bind(pelatihan); err != nil {
		return err
	}

	if pelatihan.Id_pelatihan == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	Pelatihan := &models.Tbl_pelatihan{}
	Pelatihan = pelatihan.Tbl_pelatihan

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Model(&models.Tbl_pelatihan{}).UpdateColumns(&Pelatihan).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary DeletePelatihan
@Tags Pelatihan-Controller
@Accept  json
@Produce  json
@Param id path int true "Delete Pelatihan by id"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /pelatihan/{id} [post]*/
func DeletePelatihan(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	pelatihan := &models.Tbl_pelatihan{}
	pelatihan.Id_pelatihan = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&pelatihan).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}