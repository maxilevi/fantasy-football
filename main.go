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

// @host localhost:8080
// @BasePath /api/

// @securityDefinitions.bearer BearerAuth
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
