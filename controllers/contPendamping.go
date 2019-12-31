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
	 "strconv"
	"bufio"
	"encoding/base64"	
	"io/ioutil"	
)

func GetPendamping(c echo.Context) error {
	qk 		:= c.QueryParam("id")

	var tmpPath, urlPath, blobFile, flag, host string
	flag = "PENDAMPING"
	host = c.Request().Host

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	Pendamping := models.Pendamping{}
	q := con
	q = q.Table("tbl_pendamping t1")
	q = q.Joins("join tbl_user t2 on t2.id_user = t1.id_pendamping")
	q = q.Joins("join tbl_account t3 on t3.id_user = t1.id_pendamping")
	q = q.Joins("join tbl_roles t4 on t4.id = t3.id_roles")
	q = q.Select("t2.*, t1.jenis_pendamping, t1.periode, t3.id_roles, t3.username, t3.password, t4.roles_name")
	q = q.Where("t1.id_pendamping = ?", qk)
	if ErrNo := q.Scan(&Pendamping); ErrNo.Error != nil { 
		log.Println("Erro : ", ErrNo.Error)
		return echo.ErrNotFound
	}

    // get photo user
    var photo []models.Tbl_user_files
	if err := con.Table("tbl_user_files").Where(&models.Tbl_user_files{Id_user: Pendamping.Id_user}).Where("type = ?", "IMAGE").Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

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

	Pendamping.Photo = photo

	// Pendamping.Password = s.Repeat("*", len(Pendamping.Password))
	Pendamping.Password = "******"

	r := &models.Jn{Msg: Pendamping}
	defer con.Close()

	return c.JSON(http.StatusOK, r)

}

func AddPendamping(c echo.Context) (err error) {
	Pendamping := &models.Pendamping{}

	if err := c.Bind(Pendamping); err != nil {
		return err
	}

	// get log post
	helpers.FetchPost(Pendamping)	

	// validation
	if Pendamping.Nik == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") }
	if Pendamping.Jenis_pendamping == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill jenis_pendamping") }
	if Pendamping.Periode == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill periode") }
	if Pendamping.Username == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill username") }
	if Pendamping.Password == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill password") }
	if Pendamping.Id_roles == 0 { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill id_roles") }

	// init DB Con
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// store user
	user := &models.Tbl_user{}
	user = Pendamping.Tbl_user

	if err := con.Create(&user).Error; err != nil { 
		log.Println(err)
		return echo.ErrInternalServerError 
	}

	// store pendamping
	pendamping 					:= &models.Tbl_pendamping{}
	pendamping.Id_pendamping 	= user.Id_user
	pendamping.Jenis_pendamping = Pendamping.Jenis_pendamping
	pendamping.Periode 			= Pendamping.Periode

	if err := con.Create(&pendamping).Error; err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}

	// store account_pendamping
	account 			:= &models.Tbl_account{}
	account.Id_user  	= user.Id_user
	account.Id_roles  	= Pendamping.Id_roles
	account.Username  	= Pendamping.Username
	// bycrupt password
	pwd := []byte(Pendamping.Password)
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
	account.Password  	= string(hash)

	if err := con.Create(&account).Error; err != nil {
		log.Println(err)
		return echo.ErrInternalServerError
	}		

	// close DB Con
	defer con.Close()

	r := &models.Jn1{Msg: "Success Store Data", Id: user.Id_user}
	return c.JSON(http.StatusOK, r)	
}

func UpdatePendamping(c echo.Context) (err error) {
	Pendamping := &models.Pendamping{}

	if err := c.Bind(Pendamping); err != nil {
		return err
	}

	if Pendamping.Id_user == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Please, fill id")
	}

	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// update user
	user := &models.Tbl_user{}
	user = Pendamping.Tbl_user

	if err := con.Model(&models.Tbl_user{}).UpdateColumns(&user).Error; err != nil {
		return echo.ErrInternalServerError
	}

	// update pendamping
	pendamping 					:= &models.Tbl_pendamping{}
	pendamping.Id_pendamping 	= user.Id_user
	pendamping.Jenis_pendamping = Pendamping.Jenis_pendamping
	pendamping.Periode 			= Pendamping.Periode

	if err := con.Model(&models.Tbl_pendamping{}).UpdateColumns(&pendamping).Error; err != nil {
		return echo.ErrInternalServerError
	}

	// update account
	account 			:= &models.Tbl_account{}
	account.Id_user  	= user.Id_user
	account.Id_roles  	= Pendamping.Id_roles
	account.Username  	= Pendamping.Username
	// bycrupt password
	pwd := []byte(Pendamping.Password)
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

func GetPaginatePendamping(c echo.Context) (err error) {	
	if err := helpers.PaginatePendamping(c, &r); err != nil {
		return err
	}	
	return c.JSON(http.StatusOK, r)
}

func UploadPendampingFiles(c echo.Context) (err error) {
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
		PendampingFiles := &models.Tbl_user_files{}

		PendampingFiles.Id_user = id
		PendampingFiles.Files   = encoded
		PendampingFiles.Description = description
		PendampingFiles.Type = types
		PendampingFiles.Is_display = is_display

		if err := con.Create(&PendampingFiles).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
	}

	defer con.Close()

	log.Println("Uploads Pendamping's file to id : ", id)
	r := &models.Jn{Msg: "Success Upload files"}	
	return c.JSON(http.StatusOK, r)	
	
}