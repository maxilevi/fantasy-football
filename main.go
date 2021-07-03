package main

import (
	"./app"
	"github.com/joho/godotenv"
	"log"
	"os"
)

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
