package controllers

import (
	"net/http"
	"github.com/labstack/echo"

	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/models"
	"uepkube-api/db"
	"uepkube-api/helpers"
	"strconv"
	"log"
	"fmt"

	"bufio"
	"encoding/base64"	
	"io/ioutil"		
)

type Tbl_pendamping struct {
	*models.Tbl_pendamping
	Nama string `json:"nama"`
}

type Tbl_periods_uepkube struct {
	*models.Tbl_periods_uepkube
	BantuanPeriods *Tbl_bantuan_periods `json:"bantuan_periods" gorm:"foreignkey:id;association_foreignkey:id_periods"`
}

type Tbl_bantuan_periods struct {
	*models.Tbl_bantuan_periods
	CreditDebit 	[]*models.Tbl_credit_debit `json:"credit_debit" gorm:"foreignkey:id_periods;association_foreignkey:id"`
}

type Tbl_uep struct {
	*models.Tbl_uep
	Pendamping 		*Tbl_pendamping `json:"pendamping" gorm:"foreignkey:id_pendamping;association_foreignkey:id_pendamping"`
	JenisUsaha 		*models.Tbl_jenis_usaha `json:"jenis_usaha" gorm:"foreignkey:id_jenis_usaha;association_foreignkey:id_usaha"`
	PeriodsHistory 	[]*Tbl_periods_uepkube `json:"periods_history" gorm:"foreignkey:id_uep"`	
	Photo 			[]*models.Tbl_uepkube_photo `json:"photo" gorm:"foreignkey:id_uep"`
}

type Tbl_user struct {
	*models.Tbl_user
	*Tbl_uep
	Kelurahan 	*models.Tbl_kelurahan `json:"kelurahan" gorm:"foreignkey:id_kelurahan;association_foreignkey:id_kelurahan"`
	Kecamatan 	*models.Tbl_kecamatan `json:"kecamatan" gorm:"foreignkey:id_kecamatan;association_foreignkey:id_kecamatan"`
	Kabupaten 	*models.Tbl_kabupaten `json:"kabupaten" gorm:"foreignkey:id_kabupaten;association_foreignkey:id_kabupaten"`
}

/*@Summary GetUepById
@Tags Uep-Controller
@Accept  json
@Produce  json
@Param id query int true "int"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /uep [get]*/
func GetUep(c echo.Context) error {
	id 		:= c.QueryParam("id")

	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	
	
	User 	:= Tbl_user{}
	q := con
	q = q.Model(&User)
	q = q.Joins("join tbl_uep on tbl_uep.id_uep = tbl_user.id_user")
	q = q.Select("tbl_uep.*, tbl_user.*")
	q = q.Preload("JenisUsaha")
	q = q.Preload("PeriodsHistory.BantuanPeriods.CreditDebit", func(q *gorm.DB) *gorm.DB {
		return q.Where("id_uep = ?", id)
	})
	q = q.Preload("Pendamping")
	q = q.Preload("Kelurahan")
	q = q.Preload("Kecamatan")
	q = q.Preload("Kabupaten")
	q = q.Preload("Photo", func(q *gorm.DB) *gorm.DB {
		return q.Where("id_uep = ?", id)	
	})
	q = q.First(&User, id)

	for i,_ := range User.Photo {
			ImageBlob := User.Photo[i].Photo
			User.Photo[i].Photo = "data:image/png;base64," + ImageBlob	
		}

	//get all transac credit_debit
	// con.Table("tbl_credit_debit").Select("*").Where("id_uep = ?", id).Scan(&User.BantuanPeriods.CreditDebit)

	r := &models.Jn{Msg: User}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

/*@Summary GetPaginateUep
@Tags Uep-Controller
@Accept  json
@Produce  json
@Param uep body models.PosPagin true "Show Uep Paginate"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /uep [post]*/
func GetPaginateUep(c echo.Context) (err error) {	
	if err := helpers.PaginateUep(c, &r); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

/*@Summary AddUep
@Tags Uep-Controller
@Accept  json
@Produce  json
@Param uep body models.Uep true "Add Uep"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /uep/add [post]*/
func AddUep(c echo.Context) (err error) {
	Uep := &models.Uep{}

	if err := c.Bind(Uep); err != nil {
		return err
	}

	// log.Println("Uploaded : ", Uep.Photo)

	user := &models.Tbl_user{}
	user = Uep.Tbl_user

	// validation
	if Uep.Id_pendamping == 0 { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id Pendamping") }
	if Uep.Id_jenis_usaha == 0 { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Id jenis usaha") }
	if Uep.Id_periods == 0 { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Bantuan Modal") }
	if Uep.Nama_usaha == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Nama Usaha Modal") }	
	if Uep.Nik == "" { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") 
	} else {
		if len(Uep.Nik) > 16 || len(Uep.Nik) < 16 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please fill NIK with 16 Digits")
		}
	}

	uep := &models.Tbl_uep{}
	uep.Id_pendamping = Uep.Id_pendamping
	uep.Nama_usaha = Uep.Nama_usaha
	uep.Id_jenis_usaha = Uep.Id_jenis_usaha
	uep.Status = Uep.Status

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// check if nik is not exist
    var nik []string
	con.Table("tbl_user").Where("nik = ?", Uep.Nik).Pluck("nik", &nik)
	if len(nik) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maaf NIK sudah digunakan")
	}

	// create user
	if err := con.Create(&user).Error; err != nil {return echo.ErrInternalServerError}

	// create uep
	uep.Id_uep = user.Id_user
	if err := con.Create(&uep).Error; err != nil {return echo.ErrInternalServerError}

	// store bantuan_periods_history
	uepPeriods := &models.Tbl_periods_uepkube{}
	uepPeriods.Id_uep = user.Id_user
	uepPeriods.Id_periods = Uep.Id_periods
	if err := con.Create(&uepPeriods).Error; err != nil {return echo.ErrInternalServerError}

	// store creditDebit
	creditDebit := &models.Tbl_credit_debit{}
	creditDebit.Id_uep = user.Id_user
	creditDebit.Debit = 1
	creditDebit.Id_periods = Uep.Id_periods
	var nilai []float32
	con.Table("tbl_bantuan_periods").Where("id = ?", creditDebit.Id_periods).Pluck("bantuan_modal", &nilai)
	creditDebit.Nilai = nilai[0]
	creditDebit.Deskripsi = fmt.Sprintf("Credit dengan nilai : Rp. %.2f,-", nilai[0])
	if err := con.Create(&creditDebit).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary UpdateUep
@Tags Uep-Controller
@Accept  json
@Produce  json
@Param uep body models.Uep true "Update Uep"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /uep [put]*/
func UpdateUep(c echo.Context) (err error) {
	Uep := &models.Uep{}

	if err := c.Bind(Uep); err != nil {
		return err
	}

	// validation
	if Uep.Id_user == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}
	if Uep.Nik == "" { 
		return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") 
	} else {
		if len(Uep.Nik) > 16 || len(Uep.Nik) < 16 {
			return echo.NewHTTPError(http.StatusBadRequest, "Please fill NIK with 16 Digits")
		}
	}	
	
	user := &models.Tbl_user{}
	user = Uep.Tbl_user

	uep := &models.Tbl_uep{}
	uep.Id_pendamping = Uep.Id_pendamping
	uep.Nama_usaha = Uep.Nama_usaha
	uep.Id_jenis_usaha = Uep.Id_jenis_usaha
	uep.Status = Uep.Status

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// update user
	if err := con.Model(&models.Tbl_user{}).UpdateColumns(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	// update uep
	uep.Id_uep = user.Id_user
	if err := con.Model(&models.Tbl_uep{}).UpdateColumns(&uep).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

/*@Summary DeleteUep
@Tags Uep-Controller
@Accept  json
@Produce  json
@Param id path int true "Delete Uep by id"
@Success 200 {object} models.Jn
@Failure 400 {object} models.HTTPError
@Failure 401 {object} models.HTTPError
@Failure 404 {object} models.HTTPError
@Failure 500 {object} models.HTTPError
@security ApiKeyAuth
@Router /uep/{id} [post]*/
func DeleteUep(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	user := &models.Tbl_user{}
	user.Id_user = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// delete user
	if err := con.Delete(&user).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
	
	// delete user_photo
	if err := con.Where("id_user = ?", user.Id_user).Delete(models.Tbl_user_photo{}).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data "	}
	return c.JSON(http.StatusOK, r)	
}

func UploadUepFiles(c echo.Context) (err error) {
	id, _ 			:= strconv.Atoi(c.QueryParam("id"))
	is_display, _ 	:= strconv.Atoi(c.QueryParam("is_display"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["photo"]

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	log.Println("files : ", len(files))

	for _,f := range files {

		src, err := f.Open()
		if err != nil {
			return err
		}
		defer src.Close()	

	    // Read entire JPG into byte slice.
	    reader := bufio.NewReader(src)
	    content, _ := ioutil.ReadAll(reader)

	    // Encode as base64.
	    encoded := base64.StdEncoding.EncodeToString(content)
		
		// init photo model
		photo := &models.Tbl_user_photo{}

		photo.Id_user = id
		photo.Is_display = is_display
		photo.Photo   = encoded

		if err := con.Create(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
		
		// log.Println("encoded : ", encoded)	

	}

	defer con.Close()

	log.Println("Uploads Uep's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}