package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"uepkube-api/db"
	"uepkube-api/models"
	"uepkube-api/helpers"
	"strconv"
	// "fmt"
	"log"
	"bufio"
	"encoding/base64"	
	"io/ioutil"	
)

func GetLapKeu(c echo.Context) error {
	id 		:= c.QueryParam("id")

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Lapkeu := models.Tbl_lapkeu_uepkube{}
	q := con
	q = q.Model(&Lapkeu)
	q = q.First(&Lapkeu, id)

	r := &models.Jn{Msg: Lapkeu}
	defer con.Close()

	return c.JSON(http.StatusOK, r)
}

func GetPaginateLapKeu(c echo.Context) (err error) {	
	if err := helpers.PaginateLapkeu(c, &r); err != nil {
		return echo.ErrInternalServerError
	}	
	return c.JSON(http.StatusOK, r)
}

func AddLapKeu(c echo.Context) (err error) {
	lapkeu := &models.Lapkeu{}
	
	if err := c.Bind(lapkeu); err != nil {
		return err
	}

	// validation
	if lapkeu.Id_uep == 0 && lapkeu.Id_kube == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_uep or id_kube")
	}	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Lapkeu := &models.Tbl_lapkeu_uepkube{}
	Lapkeu = lapkeu.Tbl_lapkeu_uepkube

	if err := con.Create(&Lapkeu).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UpdateLapKeu(c echo.Context) (err error) {
	lapkeu := &models.Lapkeu{}

	if err := c.Bind(lapkeu); err != nil {
		return err
	}

	if lapkeu.Id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	Lapkeu := &models.Tbl_lapkeu_uepkube{}
	Lapkeu = lapkeu.Tbl_lapkeu_uepkube

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Model(&models.Tbl_lapkeu_uepkube{}).UpdateColumns(&Lapkeu).Error; err != nil {
		return echo.ErrInternalServerError
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

func DeleteLapKeu(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))

	if id == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "please, fill id")
	}

	lapkeu := &models.Tbl_lapkeu_uepkube{}
	lapkeu.Id = id

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	if err := con.Delete(&lapkeu).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	defer con.Close()

	r := &models.Jn{Msg: "Success Delete Data"	}
	return c.JSON(http.StatusOK, r)	
}

func UploadLapKeuFiles(c echo.Context) (err error) {
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
		LapKeuFiles := &models.Tbl_lapkeu_files{}

		LapKeuFiles.Id_lapkeu = id
		LapKeuFiles.Files   = encoded
		LapKeuFiles.Description = description
		LapKeuFiles.Type = types
		LapKeuFiles.Is_display = is_display

		if err := con.Create(&LapKeuFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
	}

	defer con.Close()

	log.Println("Uploads Lapkeu's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}