package controllers

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	 _"github.com/jinzhu/gorm/dialects/mysql"
	 "uepkube-api/models"
	 "uepkube-api/db"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ResSignin struct{
	Token string `json:"token"`
	TokenType string `json:"tokenType"`
	Id_user int `json:"id_user"`
	Nama string `json:"nama"`
	Username string `json:"usernane"`
	Group_account string `json:"roles"`
}

// @Summary SignIn
// @Tags Auth-Controller
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {object} models.Jn
// @Failure 400 {object} models.HTTPError
// @Failure 401 {object} models.HTTPError
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /auth/signin [post]
func SignIn(c echo.Context) error {
	// prepare db
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	

	// init var
	u := c.QueryParam("username")
	p := c.QueryParam("password")
	
	signin := models.Tbl_account{}
	ressignin := ResSignin{}

	// check user
	if err := con.Table("tbl_account").Where(&models.Tbl_account{Username:u}).First(&signin).Error; gorm.IsRecordNotFoundError(err) {return echo.NewHTTPError(http.StatusNotFound, "User Not Found")}

	// compare password from db
    if err := bcrypt.CompareHashAndPassword([]byte(signin.Password),[]byte(p));err != nil {return echo.NewHTTPError(http.StatusBadRequest, "Wrong Password")}

    // join+get user detail
	if err := con.Table("tbl_account").Select("tbl_user.nama, tbl_group_account.group_account, tbl_user.nama").Joins("join tbl_user on tbl_user.id_user = tbl_account.id_user").Joins("join tbl_group_account on tbl_group_account.id_group = tbl_account.id_group").Where(&models.Tbl_account{Username:u}).Scan(&ressignin).Error; gorm.IsRecordNotFoundError(err) {return echo.NewHTTPError(http.StatusNotFound, "Data Not Found")}

	// Set custom claims for UEP
	claims := &models.Claims{
		ressignin.Nama,
		ressignin.Group_account,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	ressignin.Token = t
	ressignin.TokenType = "Bearer"
	ressignin.Username = signin.Username
	ressignin.Id_user = signin.Id_user

	defer con.Close()
    return c.JSON(http.StatusOK, ressignin)
}