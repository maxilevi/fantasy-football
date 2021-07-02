package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"./models"
	"./handlers"
)

type App struct {
	address string
	router *mux.Router
	db *gorm.DB
}

func Configure() *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	handlers.AddUserRoutes(r)
	return r
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {

	}
}

func CreateApp(address string) (*App, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db")
		return nil, err
	}

	Migrate(db)

	app := App{}
	app.address = address
	app.db = db
	app.router = Configure()
	return &app, nil
}

func (a *App) Run() {
	http.Handle("/", a.router)
	log.Printf("Listening on address %s\n", a.address)
	log.Fatal(http.ListenAndServe(a.address, nil))
}

func (a *App) Close() {
	sqlDB, err := a.db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	_ = sqlDB.Close()
}