package main

import (
	"./app"
	"github.com/joho/godotenv"
	"log"
	"os"
)


// @title Fantasy football manager API
// @version 1.0
// @description Fantasy football manager microservice.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/

// @securityDefinitions.apiKey BearerAuth
// @in header
// @tokenUrl POST "/session"
// @name Authorization
// @scope.write Grants read and write access to user information
// @scope.admin Grants read and write access to administrative information
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	application, err := app.CreateApp("localhost:8080", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal("Failed to start app")
	}
	defer application.Close()
	application.Run()
}
