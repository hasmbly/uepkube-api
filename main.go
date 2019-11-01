package main

import (
	"uepkube-api/routes"
)

// @title UepKube API
// @version 1.0
// @description The Documentation for UepKube REST API, Enjoyed..

// @contact.name -> @h4sb1
// @contact.email justhasby@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api-uepkube.pusdatin-dinsos.jakarta.go.id
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := routes.Init()

	// serve on port
	e.Logger.Fatal(e.Start(":9093"))
}