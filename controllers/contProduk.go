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
	 // "log"
)

/*@Summary GetProdukById
@Tags Produk-Controller
@Accept  json
@Produce  json
@Param id query int true "int"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@Router /produk [get]*/
func GetProduk(c echo.Context) error {
	id 		:= c.QueryParam("id")
	For, _	:= strconv.Atoi(c.QueryParam("for"))

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Produks := models.ShowProduks{}

	var flag int = int(For) // flag uep : 0 | kube : 1
	var id_uep []int
	var id_kube []int

	if flag == 0 {
		q1 := con
		q1 = q1.Table("tbl_usaha_produk t1")
		q1 = q1.Where("t1.id_uep = ?", id)
		q1 = q1.Pluck("t1.id_uep", &id_uep)
		if len(id_uep) != 0 { flag = 0 }
	} else if flag == 1 {
		q2 := con
		q2 = q2.Table("tbl_usaha_produk t1")
		q2 = q2.Where("t1.id_kube = ?", id)
		q2 = q2.Pluck("t1.id_kube", &id_kube)
		if len(id_kube) != 0 { flag = 1 }
	}


	// query uep or kube
	q := con
	q = q.Table("tbl_usaha_produk t1")

	if flag == 0 {
		q = q.Select(
			"t1.id,t2.nama,t2.alamat,t2.no_hp,t3.nama_produk,t3.deskripsi,t4.jenis_usaha")
		q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_uep")
		q = q.Joins("join tbl_jenis_usaha t4 on t4.id_usaha = t1.id_usaha")
		q = q.Joins("join tbl_produk t3 on t3.id_produk = t1.id_produk")
		q = q.Joins("join tbl_produk_photo t5 on t5.id_produk = t1.id_produk")
		q = q.Where("t1.id_uep = ?", id)

		if ErrNo := q.Scan(&Produks); ErrNo.Error != nil { 
			return echo.ErrNotFound
		}

		// get All photo
		var photos []models.Tbl_produk_photo
		q3 := con
		q3 = q3.Table("tbl_usaha_produk t1")
		q3 = q3.Joins("join tbl_produk_photo t2 on t2.id_produk = t1.id_produk")
		q3 = q3.Where("t1.id_uep = ?", id)
		q3 = q3.Select("t2.id_produk, t2.id, t2.is_display, t2.photo")
		q3 = q3.Scan(&photos)

		// log.Println("photos : ", photos)

		for i,_ := range photos {
			ImageBlob := photos[i].Photo
			photos[i].Photo = "data:image/png;base64," + ImageBlob
			Produks.Photo = append(Produks.Photo, photos[i])
		}

	} else if flag == 1 {
		q = q.Select(
			"t1.id,tbl_kube.nama_kube as nama,tbl_user.alamat,tbl_user.no_hp,tbl_produk.nama_produk,tbl_produk.deskripsi,tbl_jenis_usaha.jenis_usaha,tbl_produk_photo.photo")
		q = q.Joins("join tbl_kube on tbl_kube.id_kube = t1.id_kube")
		q = q.Joins("join tbl_user on tbl_user.id_user = tbl_kube.ketua")
		q = q.Joins("join tbl_jenis_usaha on tbl_jenis_usaha.id_usaha = t1.id_usaha")
		q = q.Joins("join tbl_produk on tbl_produk.id_produk = t1.id_produk")
		q = q.Joins("join tbl_produk_photo on tbl_produk_photo.id_produk = t1.id_produk")
		q = q.Where("t1.id_kube = ?", id)

		if ErrNo := q.Scan(&Produks); ErrNo.Error != nil { 
			return echo.ErrNotFound
		}

		// get All photo
		var photos []models.Tbl_produk_photo
		q3 := con
		q3 = q3.Table("tbl_usaha_produk t1")
		q3 = q3.Joins("join tbl_produk_photo t2 on t2.id_produk = t1.id_produk")
		q3 = q3.Where("t1.id_kube = ?", id)
		q3 = q3.Select("t2.id_produk, t2.id, t2.is_display, t2.photo")
		q3 = q3.Scan(&photos)

		// log.Println("photos : ", photos)

		for i,_ := range photos {
			ImageBlob := photos[i].Photo
			photos[i].Photo = "data:image/png;base64," + ImageBlob
			Produks.Photo = append(Produks.Photo, photos[i])
		}		
	}

	r := &models.Jn{Msg: Produks}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

/*@Summary GetPaginateProduk
@Tags Produk-Controller
@Accept  json
@Produce  json
@Param produk body models.PosPagin true "Show Produk Paginate"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /produk [post]*/
func GetPaginateProduk(c echo.Context) (err error) {
	if err := helpers.PaginateProduk(c, &r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

/*@Summary AddProduk
@Tags Produk-Controller
@Accept  json
@Produce  json
@Param produk body models.Tbl_kube true "Add Produk"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /produk/add [post]*/
func AddProduk(c echo.Context) (err error) {
	produk := &models.Produk{}

	if err := c.Bind(produk); err != nil {
		return err
	}

	usaha_produk := &models.Tbl_usaha_produk{}
	usaha_produk = produk.Tbl_usaha_produk

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Create(&usaha_produk).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary UpdateProduk
@Tags Produk-Controller
@Accept  json
@Produce  json
@Param produk body models.Tbl_kube true "Update Produk"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /produk [put]*/
func UpdateProduk(c echo.Context) (err error) {
	produk := &models.Produk{}

	if err := c.Bind(produk); err != nil {
		return err
	}

	if produk.Id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	if produk.Id_uep == 0 && produk.Id_kube == 0  {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id uep or kube")
	}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// update user
	usaha_produk := &models.Tbl_usaha_produk{}
	usaha_produk = produk.Tbl_usaha_produk

	if err := con.Model(&models.Tbl_usaha_produk{}).UpdateColumns(&usaha_produk).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}
/*@Summary DeleteProduk
@Tags Produk-Controller
@Accept  json
@Produce  json
@Param id path int true "Delete Produk by id"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /produk/{id} [post]*/
func DeleteProduk(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	produk := &models.Tbl_usaha_produk{}
	produk.Id = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&produk).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}

func UploadProdukFiles(c echo.Context) (err error) {
	return nil
}

