package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"	
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/db"
	 "uepkube-api/models"
	 "uepkube-api/helpers"
	 "log"
	 "fmt"
	"bufio"
	"encoding/base64"	
	"io/ioutil"	
	 "strconv"
)

func GetVerifikator(c echo.Context) error {
	qk 		:= c.QueryParam("id")

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "VERIFIKATOR"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Verifikator := models.Verifikator{}
	q := con
	q = q.Table("tbl_account t1")
	q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_user")
	q = q.Joins("join tbl_roles t4 on t4.id = t1.id_roles")
	q = q.Select("t2.*, t1.id_roles, t1.username, t1.password, t4.roles_name")
	q = q.Where("t1.id_user = ?", qk)
	q = q.Where("t1.id_roles = ?", 2)
	if ErrNo := q.Scan(&Verifikator); ErrNo.Error != nil { 
		log.Println("Erro : ", ErrNo.Error)
		return echo.ErrNotFound
	}

    // get photo user
    var photo []models.Tbl_user_files
	if err := con.Table("tbl_user_files").Where(&models.Tbl_user_files{Id_user: Verifikator.Id_user}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			id_photo := photo[i].Id

			tmpPath	= fmt.Sprintf(helpers.GoPath + "/src/uepkube-api/static/assets/images/%s_id_%s_photo_id_%d.png", flag,qk,id_photo)
			urlPath	= fmt.Sprintf("http://%s/images/%s_id_%s_photo_id_%d.png", host,flag,qk,id_photo)
			blobFile = photo[i].Files

			if check := CreateFile(tmpPath, blobFile); check == false {
				log.Println("blob is empty : ", check)
			}
		
			photo[i].Files = urlPath

		}

	Verifikator.Photo = photo

	Verifikator.Password = "******"

	r := &models.Jn{Msg: Verifikator}
	defer con.Close()
	
	return c.JSON(http.StatusOK, r)

}

func AddVerifikator(c echo.Context) (err error) {
	Verifikator := &models.Verifikator{}

	if err := c.Bind(Verifikator); err != nil {
		return err
	}

	// get log post
	helpers.FetchPost(Verifikator)

	// validation
	if Verifikator.Nik == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") }

	// init DB Con
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// store user
	user := &models.Tbl_user{}
	user = Verifikator.Tbl_user

	if err := con.Create(&user).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// store account_verifikator
	account 			:= &models.Tbl_account{}
	account.Id_user  	= user.Id_user
	account.Id_roles  	= Verifikator.Id_roles
	account.Username  	= Verifikator.Username
	// bycrupt password
	pwd := []byte(Verifikator.Password)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }	
	account.Password  	= string(hash)

	if err := con.Create(&account).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}		
	// close DB Con
	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
	return c.JSON(http.StatusOK, r)
}

func UpdateVerifikator(c echo.Context) (err error) {
	Verifikator := &models.Verifikator{}

	if err := c.Bind(Verifikator); err != nil {
		return err
	}

	if Verifikator.Id_user == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// update user
	user := &models.Tbl_user{}
	user = Verifikator.Tbl_user

	if err := con.Model(&models.Tbl_user{}).UpdateColumns(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	// update account_verifikator
	account 			:= &models.Tbl_account{}
	account.Id_user  	= user.Id_user
	account.Id_roles  	= Verifikator.Id_roles
	account.Username  	= Verifikator.Username
	// bycrupt password
	pwd := []byte(Verifikator.Password)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }	
	account.Password  	= string(hash)

	if err := con.Model(&models.Tbl_account{}).UpdateColumns(&account).Error; err != nil {
		return echo.ErrInternalServerError
	}	

	defer con.Close()
	
	r := &models.Jn{Msg: "Success Update Data"}
	return c.JSON(http.StatusOK, r)
}

func GetPaginateVerifikator(c echo.Context) (err error) {	
	if err := helpers.PaginateVerifikator(c, &r); err != nil {
		return err
	}	
	return c.JSON(http.StatusOK, r)
}

func UploadVerifikatorFiles(c echo.Context) (err error) {
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
		VerifikatorFiles := &models.Tbl_user_files{}

		VerifikatorFiles.Id_user = id
		VerifikatorFiles.Files   = encoded
		VerifikatorFiles.Description = description
		VerifikatorFiles.Type = types
		VerifikatorFiles.Is_display = is_display

		if err := con.Create(&VerifikatorFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
	}

	defer con.Close()

	log.Println("Uploads Verifikator's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}