package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"	
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/db"
	 "uepkube-api/models"
	 "log"
)

func GetVerifikator(c echo.Context) error {
	qk 		:= c.QueryParam("id")

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
    var photo []models.Tbl_user_photo
	if err := con.Table("tbl_user_photo").Where(&models.Tbl_user_photo{Id_user: Verifikator.Id_user}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			if photo[i].Photo != "" {
				ImageBlob := photo[i].Photo
				photo[i].Photo = "data:image/png;base64," + ImageBlob	
			}

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