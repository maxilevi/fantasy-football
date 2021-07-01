package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/postgres"
	"log"
	"net/http"
	"os"
)

type App struct {
	address string
	router *mux.Router
	db *gorm.DB
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	app, err := createApp("http://localhost:8080")
	if err != nil {
		log.Fatal("Failed to start app")
	}
	defer app.Close()
	app.Run()
}

func createApp(address string) (*App, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db")
		return nil, err
	}

	migrations.Migrate()
	
	app := App{}
	app.address = address
	app.db = db
	app.router = routes.Configure()
	return &app, nil
}

func (a *App) Run() {
	http.Handle("/", a.router)
	log.Fatal(http.ListenAndServe(a.address, nil))
}

func (a *App) Close() {
	sqlDB, err := a.db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	_ = sqlDB.Close()
}