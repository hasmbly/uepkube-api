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

	"bufio"
	"encoding/base64"	
	"io/ioutil"		
)

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
	/*prepare DB*/
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	
	
	var val string
	User 	:= models.Tbl_user{}
	R 		:= models.CustU{}

	/*check if query key -> "val"*/
	qk := c.QueryParams()
	for k,v := range qk {
		if k == "val" {
			return err
		} else if k == "id" {
			val = v[0]
			id,_ := strconv.Atoi(val)
			/*find kube by Nama_kube:*/
			if err := con.Where(&models.Tbl_user{Id_user:id}).First(&User).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}
			/*find uep by id + join pendaming-user*/
			id = User.Id_user
			if err := con.Table("tbl_uep").Select("tbl_uep.bantuan_modal, tbl_uep.status,tbl_uep.id_pendamping,tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_uep.id_pendamping").Where(&models.Tbl_uep{Id_uep:id}).Scan(&R).Error; gorm.IsRecordNotFoundError(err) {
				return echo.NewHTTPError(http.StatusNotFound, "Uep Not Found")
			}			
		}
	}
	
    // get photo user
    var photo []models.Tbl_user_photo
	if err := con.Table("tbl_user_photo").Where(&models.Tbl_user_photo{Id_user: User.Id_user}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			if photo[i].Photo != "" {
				ImageBlob := photo[i].Photo
				photo[i].Photo = "data:image/png;base64," + ImageBlob	
			}

		}

	R.Photo = photo
	R.Flag = "UEP"

	r := &models.Jn{Msg: models.U{User, R}}
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
	if Uep.Bantuan_modal == 0 { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill Bantuan Modal") }
	if Uep.Nik == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") }

	uep := &models.Tbl_uep{}
	uep.Id_pendamping = Uep.Id_pendamping
	uep.Bantuan_modal = Uep.Bantuan_modal
	uep.Status = Uep.Status

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Create(&user).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	uep.Id_uep = user.Id_user

	if err := con.Create(&uep).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	

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

	if Uep.Id_user == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}
	
	user := &models.Tbl_user{}
	user = Uep.Tbl_user

	uep := &models.Tbl_uep{}
	uep.Id_pendamping = Uep.Id_pendamping
	uep.Bantuan_modal = Uep.Bantuan_modal
	uep.Status = Uep.Status

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Save(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	uep.Id_uep = user.Id_user

	if err := con.Save(&uep).Error; err != nil {
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