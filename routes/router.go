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
	e.GET("/swagger-api/*", echoSwagger.WrapHandler)

	// Route::Unauthenticated-Group
	o := e.Group("/api/v1")
	o.POST("/auth/signin", controllers.SignIn)
	o.GET("/lookup/uepkube", controllers.GetUepKube)
	// produk
	o.GET("/produk", controllers.GetProduk)
	o.POST("/lookup/uepkube/produk", controllers.GetPaginateProdukUepKube)
	// pelatihan
	o.GET("/pelatihan", controllers.GetPelatihan)
	o.POST("/lookup/uepkube/pelatihan", controllers.GetPaginatePelatihanUepKube)
	// faq
	o.GET("/lookup/faq", controllers.GeAllFaq)
	// persebaran
	o.GET("/lookup/persebaran", controllers.GeAllUepKubeDetail)	

	// Route::Restricted-Group-UEP
	u := e.Group("/api/v1")
	u.Use(middleware.JWTWithConfig(config))
	u.Use(middlewares.CheckUepRoles)
	// uep
	u.GET("/uep", controllers.GetUep)
	u.POST("/uep", controllers.GetPaginateUep)
	u.PUT("/uep", controllers.UpdateUep)
	u.POST("/uep/add", controllers.AddUep)
	u.POST("/uep/:id", controllers.DeleteUep)

	// Routes::All Roles
	a := e.Group("/api/v1")
	a.Use(middleware.JWTWithConfig(config))
	a.Use(middlewares.CheckAllRoles)	
	// produk
	a.PUT("/produk", controllers.UpdateProduk)
	a.POST("/produk/add", controllers.AddProduk)
	a.POST("/produk/:id", controllers.DeleteProduk)
	// pelatihan
	a.PUT("/pelatihan", controllers.UpdatePelatihan)
	a.POST("/pelatihan/add", controllers.AddPelatihan)
	a.POST("/pelatihan/:id", controllers.DeletePelatihan)
	// inventaris
	a.GET("/inventaris", controllers.GetInventaris)
	a.POST("/inventaris", controllers.GetPaginateInventaris)	
	a.PUT("/inventaris", controllers.UpdateInventaris)
	a.POST("/inventaris/add", controllers.AddInventaris)
	a.POST("/inventaris/:id", controllers.DeleteInventaris)
	// aktivitas
	a.GET("/aktivitas", controllers.GetAktivitas)
	a.POST("/aktivitas", controllers.GetPaginateAktivitas)	
	a.PUT("/aktivitas", controllers.UpdateAktivitas)
	a.POST("/aktivitas/add", controllers.AddAktivitas)
	a.POST("/aktivitas/:id", controllers.DeleteAktivitas)		
	
	// Route::Restricted-Group-KUBE
	k := e.Group("/api/v1")
	k.Use(middleware.JWTWithConfig(config))
	k.Use(middlewares.CheckKubeRoles)	
	k.GET("/kube", controllers.GetKube)
	k.POST("/kube", controllers.GetPaginateKube)
	k.PUT("/kube", controllers.UpdateKube)
	k.POST("/kube/add", controllers.AddKube)
	k.POST("/kube/:id", controllers.DeleteKube)	

	return e
}
