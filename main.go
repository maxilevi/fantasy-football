package main

import (
	"github.com/joho/godotenv"
	"./app"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	application, err := app.CreateApp("localhost:8080")
	if err != nil {
		log.Fatal("Failed to start app")
	}
	defer application.Close()
	application.Run()
}