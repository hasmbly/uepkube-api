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

// Middleware Custom Claims JWT
var config = middleware.JWTConfig{
	Claims:     &models.Claims{},
	SigningKey: []byte("secret"),
}

func Init() *echo.Echo {
	// initialize
	fmt.Println("Running...")
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = helpers.CustomHTTPErrorHandler

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

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
	// kelurahan-kabupaten
	o.GET("/lookup/address", controllers.GeAllAddress)
	// jenis_usaha
	o.GET("/lookup/jenis_usaha", controllers.GeAllJenisUsaha)
	// pendamping
	o.GET("/lookup/pendamping", controllers.GeAllPendamping)	
	// bantuan_periods
	o.GET("/lookup/bantuan_periods", controllers.GeAllBantuanPeriods)
	// member_pelatihan
	o.GET("/lookup/member_pelatihan", controllers.GeAllMemberPelatihan)

	// Route::Restricted-Group-UEP
	// u := e.Group("/api/v1")
	// u.Use(middleware.JWTWithConfig(config))
	// u.Use(middlewares.CheckUepRoles)
	// // uep
	// // u.GET("/uep", controllers.GetUep)
	// // u.POST("/uep", controllers.GetPaginateUep)
	// // u.PUT("/uep", controllers.UpdateUep)
	// // u.POST("/uep/add", controllers.AddUep)
	// // u.POST("/uep/:id", controllers.DeleteUep)

	// Route::Restricted-Group-KUBE
	// k := e.Group("/api/v1")
	// k.Use(middleware.JWTWithConfig(config))
	// k.Use(middlewares.CheckKubeRoles)	
	// k.GET("/kube", controllers.GetKube)
	// k.POST("/kube", controllers.GetPaginateKube)
	// k.PUT("/kube", controllers.UpdateKube)
	// k.POST("/kube/add", controllers.AddKube)
	// k.POST("/kube/:id", controllers.DeleteKube)	


	// Routes::All Roles
	a := e.Group("/api/v1")
	a.Use(middleware.JWTWithConfig(config))
	a.Use(middlewares.CheckAllRoles)	
	// produk
	// a.PUT("/produk", controllers.UpdateProduk)
	// a.POST("/produk/add", controllers.AddProduk)
	// a.POST("/produk/:id", controllers.DeleteProduk)
	// // pelatihan
	// a.PUT("/pelatihan", controllers.UpdatePelatihan)
	// a.POST("/pelatihan/add", controllers.AddPelatihan)
	// a.POST("/pelatihan/:id", controllers.DeletePelatihan)
	// // inventaris
	// a.GET("/inventaris", controllers.GetInventaris)
	// a.POST("/inventaris", controllers.GetPaginateInventaris)	
	// a.PUT("/inventaris", controllers.UpdateInventaris)
	// a.POST("/inventaris/add", controllers.AddInventaris)
	// a.POST("/inventaris/:id", controllers.DeleteInventaris)
	// // aktivitas
	// a.GET("/aktivitas", controllers.GetAktivitas)
	// a.POST("/aktivitas", controllers.GetPaginateAktivitas)	
	// a.PUT("/aktivitas", controllers.UpdateAktivitas)
	// a.POST("/aktivitas/add", controllers.AddAktivitas)
	// a.POST("/aktivitas/:id", controllers.DeleteAktivitas)

	// CRUD Pendamping, UEP, KUBE, Verifikator
	a.POST("/:key", controllers.GetPaginateItems)	
	a.GET("/:key", controllers.GetItems)
	a.POST("/add/:key", controllers.AddItems)
	a.PUT("/:key", controllers.UpdateItems)
	a.POST("/:key/:id", controllers.DeleteItems)
	
	// uploads images
	a.POST("/upload/images/:key", controllers.UploadImages)

	// uploads pdf
	// a.POST("/uploads/pdf/:key", controllers.UploadFiles)
	
	return e
}
