package routes

import (
	"uepkube-api/controllers"
	"uepkube-api/helpers"
	"uepkube-api/models"
	"uepkube-api/middlewares"
	"fmt"
	"log"
	"net/http"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"
	"golang.org/x/crypto/bcrypt"
	_ "uepkube-api/docs"

	"github.com/jinzhu/gorm"
	"uepkube-api/db"

	 // "image"
	 // "image/png"
	 // "os"
	 // "bytes"

)

func Home(c echo.Context) error {
  return c.HTML(http.StatusOK, "<a href='https://echo.labstack.com'><img height='80' src='https://cdn.labstack.com/images/echo-logo.svg'></a><br /><pre><strong>Echo</strong> v4.1.11High performance, minimalist Go web framework</pre>")}

func BycriptPass(c echo.Context) error {
	pwd := []byte(c.Param("pass"))
    hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
	return c.JSON(http.StatusOK, string(hash))
}

func GetPhoto(c echo.Context) error {
	con, err := db.CreateCon()
	if err != nil { return echo.ErrInternalServerError }
	con.SingularTable(true)	

    var imgByte []string
	if err := con.Table("tbl_user_photo").Where("id_user = ?", 6).Pluck("photo", &imgByte).Error; gorm.IsRecordNotFoundError(err) {return echo.ErrNotFound}	
   
	defer con.Close()	

	var photos []string
	for i,_ := range imgByte {
		var html string
		html = "<img src='data:image/png;base64," + imgByte[i] +
			"alt='testing img' /><br />"
		photos = append(photos, html)
	}
	return c.HTML(http.StatusOK, photos[0] + photos[1])
}

// func Shoot() (err error) {
// 	Uep := []models.Tbl_uep{}
// 	Kube := []models.Tbl_kube{}

// 	con, err := db.CreateCon()
// 	if err != nil { return echo.ErrInternalServerError }
// 	con.SingularTable(true)

// 	// uep
// 	con.Model(&Uep).Find(&Uep)
// 	if len(Uep) != 0 {
// 		for i, _ := range Uep {
// 			// Monev
// 			Periods := models.Tbl_periods_uepkube{}
// 			Periods.Id_uep = Uep[i].Id_uep
// 			Periods.Id_periods = 2
// 			if err := con.Create(&Periods).Error; err != nil {return echo.ErrInternalServerError}	
// 		}
// 	}

// 	// kube
// 	con.Model(&Kube).Find(&Kube)
// 	if len(Kube) != 0 {
// 		for i, _ := range Kube {
// 			// Monev
// 			Periods := models.Tbl_periods_uepkube{}
// 			Periods.Id_kube = Kube[i].Id_kube
// 			Periods.Id_periods = 1
// 			if err := con.Create(&Periods).Error; err != nil {return echo.ErrInternalServerError}	
// 		}
// 	}	
// 	defer con.Close()
// 	return err
// }

// Middleware Custom Claims JWT
var config = middleware.JWTConfig{
	Claims:     &models.Claims{},
	SigningKey: []byte("secret"),
}

func Init() *echo.Echo {
	// initialize
	fmt.Println("Running...")
	e := echo.New()

	// GOPATH
	goPath := helpers.GoPath
	log.Println("gopath : ", goPath)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = helpers.CustomHTTPErrorHandler

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// for pdf
	e.Static("/pdf", goPath + "/src/uepkube-api/static/assets/pdf")
	// for pdf
	e.Static("/images", goPath + "/src/uepkube-api/static/assets/images")

	e.GET("/", Home)
	e.GET("/bycrypt/:pass", BycriptPass)
	e.GET("/photo", GetPhoto)
	e.GET("/swagger-api/*", echoSwagger.WrapHandler)

	// Route::Unauthenticated-Group
	o := e.Group("/api/v1")
	o.POST("/auth/signin", controllers.SignIn)
	o.GET("/lookup/uepkube", controllers.GetUepKube)
	// produk
	// o.GET("/produk", controllers.GetProduk)
	o.POST("/lookup/uepkube/produk", controllers.GetPaginateProdukUepKube)
	// pelatihan
	// o.GET("/pelatihan", controllers.GetPelatihan)
	o.POST("/lookup/uepkube/pelatihan", controllers.GetPaginatePelatihanUepKube)
	// faq
	o.GET("/lookup/faq", controllers.GeAllFaq)
	// persebaran
	o.GET("/lookup/persebaran", controllers.GeAllUepKubeDetail)
	// address -> kelurahan-kabupaten
	o.GET("/lookup/address", controllers.GeAllAddress)
	// address_detail
	o.GET("/lookup/address_detail", controllers.GetAllAddressDetail)
	// jenis_usaha
	o.GET("/lookup/jenis_usaha", controllers.GeAllJenisUsaha)
	// pendamping
	o.GET("/lookup/pendamping", controllers.GeAllPendamping)	
	// bantuan_periods
	o.GET("/lookup/bantuan_periods", controllers.GeAllBantuanPeriods)
	// member_pelatihan
	o.GET("/lookup/member_pelatihan", controllers.GeAllMemberPelatihan)
	// monev_pertanyaan+score indikator
	o.GET("/lookup/monev_items", controllers.GeAllMonevItems)
	// users
	o.GET("/lookup/users", controllers.GeAllUser)
	// chart_dashboard
	o.GET("/lookup/chart_dashboard", controllers.GetChartDasboard)	

	// Routes::All Roles
	a := e.Group("/api/v1")
	a.Use(middleware.JWTWithConfig(config))
	a.Use(middlewares.CheckAllRoles)

	// CRUD Universal
	a.POST("/:key", controllers.GetPaginateItems)
	a.GET("/:key", controllers.GetItems)
	a.POST("/add/:key", controllers.AddItems)
	a.PUT("/:key", controllers.UpdateItems)
	a.POST("/:key/:id", controllers.DeleteItems)
	
	// uploads files
	a.POST("/upload/files/:key", controllers.UploadFiles)
	// download files
	a.GET("/download/files/:key", controllers.DownloadFiles)

	return e
}
