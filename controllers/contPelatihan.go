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
	 "fmt"
	 // "io"
	 // "os"

	"bufio"
	"encoding/base64"	
	"io/ioutil"	
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
	// helpers.GoPath := os.Getenv("GOPATH")

	var tmpPath, urlPath, blobFile,flag,host string
	flag = "PELATIHAN"
	host = c.Request().Host	

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Pelatihan := models.Tbl_pelatihan{}
	q := con
	q = q.Model(&Pelatihan)
	q = q.Preload("Files", func(q *gorm.DB) *gorm.DB {
		return q.Where("id_pelatihan = ?", id).Where("type = 'PDF' ")
	})
	q = q.Preload("Photo", func(q *gorm.DB) *gorm.DB {
		return q.Where("id_pelatihan = ?", id).Where("type = 'IMAGE' ")	
	})	
	q = q.First(&Pelatihan, id)

	for i,_ := range Pelatihan.Photo {

			if Pelatihan.Photo[i].Files != "" {
				Pelatihan.Photo[i].Files = "data:image/png;base64," + Pelatihan.Photo[i].Files
			}
		}
	
	// get Kehadiran Pelatihan
	if err := con.Table("tbl_kehadiran t1").Select("t1.*, t2.nama").Joins("join tbl_user t2 on t2.id_user = t1.id_user").Where("t1.id_pelatihan = ?", Pelatihan.Id_pelatihan).Scan(&Pelatihan.Kehadiran).Error; err != nil { return echo.ErrInternalServerError }

	// docs pdf
	for i, _ := range Pelatihan.Files {
		id_pdf := Pelatihan.Files[i].Id
		
		tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/pdf/%s_id_%s_pdf_id_%d.pdf", flag,id,id_pdf)
		urlPath	= fmt.Sprintf("http://%s/pdf/%s_id_%s_pdf_id_%d.pdf", host,flag,id,id_pdf)
		blobFile = Pelatihan.Files[i].Files

		if check := CreateFile(tmpPath, blobFile); check == false {
			log.Println("blob is empty : ", check)
		}

		Pelatihan.Files[i].Files = urlPath
	}
	
	// photos
	// for i, _ := range Pelatihan.Files {
	// 	tmpPath := fmt.Sprintf("static/assets/pdf/%d_pelatihan.pdf", i)
	// 	urlPath := "http://" + c.Request().Host + "/pdf/" + strconv.Itoa(i) + "_pelatihan.pdf"
	// 	blobFile := Pelatihan.Files[i].Files

	// 	if check := CreateFile(tmpPath, blobFile); check == false {
	// 		log.Println("blob is empty : ", check)
	// 	}

	// 	Pelatihan.Files[i].Files = urlPath
	// }	

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

func AddPelatihanKehadiran(c echo.Context) (err error) {
	Kehadiran := &models.Tbl_kehadiran{}
	kehadirans := &models.Kehadiran{}

	if err := c.Bind(kehadirans); err != nil {
		return err
	}

	if kehadirans.Id_pelatihan == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_pelatihan")
	}

	if kehadirans.Id_pendamping == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id_pendamping")
	}	

	// id pelatihan
	id_pelatihan := kehadirans.Id_pelatihan
	id_pendamping := kehadirans.Id_pendamping

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// check if actived or not
	var status []int
	if err := con.Table("tbl_pelatihan").Where("id_pelatihan = ?", id_pelatihan).Pluck("status", &status).Error; err != nil { return echo.ErrInternalServerError }	
	if status[0] == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Maaf Pelatihan sudah Tidak Aktif")
	}

	// check if pelatihan is already started

	// check if user is already added
	for i, _ := range kehadirans.Kehadiran {
		var id []int
		var nama []string
		if err := con.Table("tbl_kehadiran").Where("id_user = ?", kehadirans.Kehadiran[i].Id_user).Pluck("id_user", &id).Error; err != nil { return echo.ErrInternalServerError }
		if len(id) != 0 {
			con.Table("tbl_user").Where("id_user = ?", id[0]).Pluck("nama", &nama)
			return echo.NewHTTPError(http.StatusBadRequest, "Maaf User " + nama[0] + " telah terdaftar dalam pelatihan ini")
		}
	}

	// do if quota && started is FALSE
	for i, _ := range kehadirans.Kehadiran {
		kehadirans.Kehadiran[i].Id_pelatihan = id_pelatihan
		kehadirans.Kehadiran[i].Id_pendamping = id_pendamping

		// check if user input same id for uep and kube
		count := i
		for ii := count+1; ii < len(kehadirans.Kehadiran); ii++ {
			if kehadirans.Kehadiran[i].Id_user == kehadirans.Kehadiran[ii].Id_user {
				return echo.NewHTTPError(http.StatusBadRequest, "Maaf Ada Id User yang sama")
			}
		}

		// check quota
		var quota []int
		if err := con.Table("tbl_pelatihan").Where("id_pelatihan = ?", kehadirans.Kehadiran[i].Id_pelatihan).Pluck("quota", &quota).Error; err != nil { return echo.ErrInternalServerError }		

		// check if quota is empty set status to not actived
		if quota[0] == 0 {
			// set status pelatihan = 0 (NOT ACTIVE)
			// if err := con.Table("tbl_pelatihan").Where("id_pelatihan = ?", id).UpdateColumn("status", 0).Error; err != nil { return echo.ErrInternalServerError }
			return echo.NewHTTPError(http.StatusBadRequest, "Kuota Pelatihan telah habis")
		}

		// check if mmember kehadiran is bigger than quota
		if len(kehadirans.Kehadiran) > len(quota) {
			return echo.NewHTTPError(http.StatusBadRequest, "Maaf Kapasitas Kehadiran User melebihi Kuota pelatihan")
		}

		// store kehadiran pelatihan
		Kehadiran = &kehadirans.Kehadiran[i]
		log.Println(Kehadiran)
		if err := con.Create(&Kehadiran).Error; err != nil { return echo.ErrInternalServerError }

		// decrease quota after adding user
		if len(quota) != 0 {
			if quota[0] != 0 {
				if err := con.Table("tbl_pelatihan").Where("id_pelatihan = ?", kehadirans.Kehadiran[i].Id_pelatihan).UpdateColumn("quota", quota[0]-1).Error; err != nil { return echo.ErrInternalServerError }			
			}
		}
	}

	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UploadPelatihanFiles(c echo.Context) (err error) {
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
		PelatihanFiles := &models.Tbl_pelatihan_files{}

		PelatihanFiles.Id_pelatihan = id
		PelatihanFiles.Files   = encoded
		PelatihanFiles.Description = description
		PelatihanFiles.Type = types
		PelatihanFiles.Is_display = is_display

		if err := con.Create(&PelatihanFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
		
		// log.Println("encoded : ", encoded)	

	}

	defer con.Close()

	log.Println("Uploads Pelatihan's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}

func DownloadPelatihanFiles(c echo.Context) (err error) {
	id 		:= c.QueryParam("id")
	
	var tmpPath, urlPath, blobFile,flag,host string
	flag = "PELATIHAN"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// get pdf blob
	PelatihanFile := []models.Tbl_pelatihan_files{}
	q := con
	q = q.Table("tbl_pelatihan_files")
	q = q.Where("id_pelatihan = ?", id)
	q = q.Find(&PelatihanFile)
	
	for i, _ := range PelatihanFile {

		if PelatihanFile[i].Type == "PDF" {

			id_pdf := PelatihanFile[i].Id
			
			tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/pdf/%s_id_%s_pdf_id_%d.pdf", flag,id,id_pdf)
			urlPath	= fmt.Sprintf("http://%s/pdf/%s_id_%s_pdf_id_%d.pdf", host,flag,id,id_pdf)
			blobFile = PelatihanFile[i].Files

			if check := CreateFile(tmpPath, blobFile); check == false {
				log.Println("blob is empty : ", check)
			}

			PelatihanFile[i].Files = urlPath

		} else if PelatihanFile[i].Type == "IMAGE" {

			id_photo := PelatihanFile[i].Id
			
			tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%s_photo_id_%d.png", flag,id,id_photo)
			urlPath	= fmt.Sprintf("http://%s/images/%s_id_%s_photo_id_%d.png", host,flag,id,id_photo)
			blobFile = PelatihanFile[i].Files


			if check := CreateFile(tmpPath, blobFile); check == false {
				log.Println("blob is empty : ", check)
			}
		
			PelatihanFile[i].Files = urlPath
		}
	}

	defer con.Close()

	r := &models.Jn{Msg: PelatihanFile}	
	return c.JSON(http.StatusOK, r)
	// return nil			
}