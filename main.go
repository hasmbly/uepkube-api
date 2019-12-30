package main

import (
	"uepkube-api/routes"
	"log"
	"os"
)

// @title UepKube API
// @version 1.0
// @description The Documentation for UepKube REST API, Enjoyed..

// @contact.name -> @h4sb1
// @contact.email justhasby@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 157.245.55.185:9093
/*@host localhost:9093*/

// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	e := routes.Init()

	// log
	f, err := os.OpenFile("err.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    e.Logger.Fatal("error opening file: %v", err)
	}
	defer f.Close()

	e.Logger.SetOutput(f)

	e.Logger.SetHeader("${time_rfc3339} ${level}")
	
	log.Println("Go Started...")

	// serve on port 
	e.Logger.Fatal(e.Start(":9093"))
}	