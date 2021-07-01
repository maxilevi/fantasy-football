package app

import (
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type App struct {
	Db *DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app, err := createApp()
	if err != nil {
		log.Fatal("Failed to start app")
	}
	app.run()
}

func createApp() (*App, error) {
	app := App{}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db")
		return nil, err
	}
	app.Db = db
	return &app, nil
}

func (a *App) run() {
	http.Handle("/", routes.Router())
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}