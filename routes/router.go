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

	// "github.com/jinzhu/gorm"
	// "uepkube-api/db"

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

// func monevDummy(c echo.Context) error {
// 	con, err := db.CreateCon()
// 	if err != nil { return echo.ErrInternalServerError }
// 	con.SingularTable(true)	

// 	Uep 	:= []models.Tbl_uep{}
// 	Kube 	:= []models.Tbl_kube{}

// 	// get Uep
// 	if err := con.Find(&Uep).Error; err != nil { 
// 		return echo.NewHTTPError(http.StatusInternalServerError, err) 
// 	} else if err == nil {
// 		for i, _ := range Uep {
// 			monevUep := &models.Tbl_monev_final{}
// 			monevUep.Id_uep 		= Uep[i].Id_uep
// 			monevUep.Id_pendamping 	= Uep[i].Id_pendamping
// 			monevUep.Is_monev 		= "BELUM"
// 			monevUep.Flag 			= "UEP"
// 			if err := con.Create(&monevUep).Error; err != nil {return echo.ErrInternalServerError}
// 		}
// 	}

// 	if err := con.Find(&Kube).Error; err != nil { 
// 		return echo.NewHTTPError(http.StatusInternalServerError, err) 
// 	} else if err == nil {
// 		for i, _ := range Kube {
// 			monevKube := &models.Tbl_monev_final{}
// 			monevKube.Id_kube 		= Kube[i].Id_kube
// 			monevKube.Id_pendamping = Kube[i].Id_pendamping
// 			monevKube.Is_monev 		= "BELUM"
// 			monevKube.Flag 			= "KUBE"
// 			if err := con.Create(&monevKube).Error; err != nil {return echo.ErrInternalServerError}
// 		}
// 	}	

// 	defer con.Close()	
// 	return c.JSON(http.StatusOK, "Success")
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
	// e.GET("/monevDummy", monevDummy)
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
	// pkt
	o.POST("/pkt", controllers.GetPaginatePkt)

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
