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

type Tbl_user struct {
	*models.Tbl_user
	*models.Tbl_uep
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
	q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha", func(q *gorm.DB) *gorm.DB {
		return q.Where("id_uep = ?", id).Preload("JenisUsaha")
	})
	q = q.Preload("PeriodsHistory.BantuanPeriods.Usaha.AllProduk.DetailProduk.JenisProduk")
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
			ImageBlob := User.Photo[i].Files
			User.Photo[i].Files = "data:image/png;base64," + ImageBlob	
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

	// log.Println("Uep_tbl_user : ", *Uep.Tbl_user)	
	// log.Println("Uep : ", Uep)	

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

	// check pendamping
    var id_pendamping []string
	con.Table("tbl_pendamping").Where("id_pendamping = ?", Uep.Id_pendamping).Pluck("id_pendamping", &id_pendamping)
	if len(id_pendamping) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maaf id_pendamping tidak ada ditemukan")
	}	

	// create user
	if err := con.Create(&user).Error; err != nil {return echo.ErrInternalServerError}

	// create uep
	uep.Id_uep = user.Id_user
	if err := con.Create(&uep).Error; err != nil {return echo.ErrInternalServerError}

	// store usaha_uep
	uepUsaha := &models.Tbl_usaha_uepkube{}
	uepUsaha.Id_uep = user.Id_user
	uepUsaha.Nama_usaha = Uep.Nama_usaha
	uepUsaha.Id_jenis_usaha = Uep.Id_jenis_usaha
	uepUsaha.Id_periods = Uep.Id_periods
	if err := con.Create(&uepUsaha).Error; err != nil {return echo.ErrInternalServerError}	

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

	// add queue monev_uepkube
	monev := &models.Tbl_monev_uepkube{}
	monev.Id_uep = user.Id_user
	monev.Id_pendamping = Uep.Id_pendamping
	monev.Is_monev = "BELUM"
	monev.Id_periods = Uep.Id_periods
	if err := con.Create(&monev).Error; err != nil {return echo.ErrInternalServerError}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data", Id: user.Id_user}
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
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id_user")
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

	// store usaha_uep
	uepUsaha := &models.Tbl_usaha_uepkube{}
	uepUsaha.Id_uep = user.Id_user
	uepUsaha.Nama_usaha = Uep.Nama_usaha
	uepUsaha.Id_jenis_usaha = Uep.Id_jenis_usaha
	uepUsaha.Id_periods = Uep.Id_periods
	if err := con.Model(&models.Tbl_usaha_uepkube{}).Where("id_uep = ?", uepUsaha.Id_uep).Where("id_jenis_usaha = ?", uepUsaha.Id_jenis_usaha).Where("id_periods = ?", uepUsaha.Id_periods).UpdateColumns(&uepUsaha).Error; err != nil {
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
	
	// delete user_UepFiles
	if err := con.Where("id_user = ?", user.Id_user).Delete(models.Tbl_user_photo{}).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data "	}
	return c.JSON(http.StatusOK, r)	
}

func UploadUepFiles(c echo.Context) (err error) {
	// query
	id, _ 			:= strconv.Atoi(c.QueryParam("id"))
	is_display, _ 	:= strconv.Atoi(c.QueryParam("is_display"))

	// formValue
	description := c.FormValue("description")
	types 		:= c.FormValue("type")

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]

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
		
		// init UepFiles model
		UepFiles := &models.Tbl_uepkube_files{}

		UepFiles.Id_uep = id
		UepFiles.Files   = encoded
		UepFiles.Description = description
		UepFiles.Type = types
		UepFiles.Is_display = is_display

		if err := con.Create(&UepFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
		
		// log.Println("encoded : ", encoded)	

	}

	defer con.Close()

	log.Println("Uploads Uep's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}