package routes

import (
	"uepkube-api/controllers"
	"uepkube-api/helpers"
	"uepkube-api/models"
	"uepkube-api/middlewares"
	"fmt"
	"log"
	//"html/template"
	//"io"
	"net/http"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"
	"golang.org/x/crypto/bcrypt"
	_ "uepkube-api/docs"
)

//type TemplateRenderer struct {
//	templates *template.Template
//}

//func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
//	return t.templates.ExecuteTemplate(w, name, data)
//}

func Home(c echo.Context) error {
  return c.HTML(http.StatusOK, "<pre><strong>Echo</strong>v4.1.11High performance, minimalist Go web framework</pre>")}

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

	// render html
//	renderer := &TemplateRenderer{
//	  templates: template.Must(template.ParseGlob("./static/views/*.html")),
//	}
//	e.Renderer = renderer

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = helpers.CustomHTTPErrorHandler

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Route::Unauthenticated-Group
	e.GET("/", Home)
	e.GET("/bycrypt/:pass", BycriptPass)
	e.GET("/swagger-api/*", echoSwagger.WrapHandler)
	e.POST("/api/v1/auth/signin", controllers.SignIn)
	e.GET("/api/v1/lookup/uepkube", controllers.GetUepOrKube)

	// Route::Restricted-Group-UEP
	u := e.Group("/api/v1")
	u.Use(middleware.JWTWithConfig(config))
	u.Use(middlewares.CheckUepRoles)
	u.GET("/uep", controllers.GetUep)
	u.POST("/uep", controllers.GetPaginateUep)
	u.POST("/uep/add", controllers.AddUep)
	
	// Route::Restricted-Group-KUBE
	k := e.Group("/api/v1")
	k.Use(middleware.JWTWithConfig(config))
	k.Use(middlewares.CheckKubeRoles)	
	k.GET("/kube", controllers.GetKube)	

	return e
}
