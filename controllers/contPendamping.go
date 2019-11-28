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

func GetPendamping(c echo.Context) error {
	qk 		:= c.QueryParam("id")

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
    var photo []models.Tbl_user_photo
	if err := con.Table("tbl_user_photo").Where(&models.Tbl_user_photo{Id_user: Pendamping.Id_user}).Find(&photo).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	for i,_ := range photo {

			if photo[i].Photo != "" {
				ImageBlob := photo[i].Photo
				photo[i].Photo = "data:image/png;base64," + ImageBlob	
			}

		}

	// Pendamping.Password = s.Repeat("*", len(Pendamping.Password))
	Pendamping.Password = "******"
	Pendamping.Photo = photo

	r := &models.Jn{Msg: Pendamping}
	defer con.Close()

	return c.JSON(http.StatusOK, r)

}

func AddPendamping(c echo.Context) (err error) {
	Pendamping := &models.Pendamping{}

	if err := c.Bind(Pendamping); err != nil {
		return err
	}

	// validation
	if Pendamping.Nik == "" { return echo.NewHTTPError(http.StatusBadRequest, "Please Fill NIK") }

	// init DB Con
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)

	// store user
	user := &models.Tbl_user{}
	user = Pendamping.Tbl_user

	if err := con.Create(&user).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

	// store pendamping
	pendamping 					:= &models.Tbl_pendamping{}
	pendamping.Id_pendamping 	= user.Id_user
	pendamping.Jenis_pendamping = Pendamping.Jenis_pendamping
	pendamping.Periode 			= Pendamping.Periode

	if err := con.Create(&pendamping).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}

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

	if err := con.Create(&account).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}		
	// close DB Con
	defer con.Close()

	r := &models.Jn{Msg: "Success Store Data"}
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