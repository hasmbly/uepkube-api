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
	 // "os"
	 "fmt"
	"bufio"
	"encoding/base64"	
	"io/ioutil"	
)

func GetAktivitas(c echo.Context) error {
	id 		:= c.QueryParam("id")

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "ACTIVITY"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Activity := models.Tbl_activity{}
	q := con
	q = q.Model(&Activity)
	q = q.Preload("Photo")
	q = q.First(&Activity, id)

	for i, _ := range Activity.Photo {
			id_photo := Activity.Photo[i].Id

			tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%s_photo_id_%d.png", flag,id,id_photo)
			urlPath	= fmt.Sprintf("http://%s/images/%s_id_%s_photo_id_%d.png", host,flag,id,id_photo)
			blobFile = Activity.Photo[i].Files

			if check := CreateFile(tmpPath, blobFile); check == false {
				log.Println("blob is empty : ", check)
			}
		
			Activity.Photo[i].Files = urlPath
	}

	r := &models.Jn{Msg: Activity}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

func GetPaginateAktivitas(c echo.Context) (err error) {	
	if err := helpers.PaginateAktivitas(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
}

func AddAktivitas(c echo.Context) (err error) {
	activity := &models.Activity{}

	if err := c.Bind(activity); err != nil {
		return err
	}

	// get log post
	helpers.FetchPost(activity)	

	Activity := &models.Tbl_activity{}
	Activity = activity.Tbl_activity

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Create(&Activity).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UpdateAktivitas(c echo.Context) (err error) {
	activity := &models.Activity{}

	if err := c.Bind(activity); err != nil {
		return err
	}

	if activity.Id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	Activity := &models.Tbl_activity{}
	Activity = activity.Tbl_activity

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Model(&models.Tbl_activity{}).UpdateColumns(&Activity).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

func DeleteAktivitas(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	activity := &models.Tbl_activity{}
	activity.Id = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&activity).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}

func UploadAktivitasFiles(c echo.Context) (err error) {
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
		
		// execute
		ActivityFiles := &models.Tbl_activity_files{}

		ActivityFiles.Id_activity = id
		ActivityFiles.Files   = encoded
		ActivityFiles.Description = description
		ActivityFiles.Type = types
		ActivityFiles.Is_display = is_display

		if err := con.Create(&ActivityFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
		
		// log.Println("encoded : ", encoded)	

	}

	defer con.Close()

	log.Println("Uploads Activity's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}