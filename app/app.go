package app

import (
	"./handlers"
	"./middleware"
	"./models"
	"./repos"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
)

type App struct {
	address   string
	router    *mux.Router
	db        *gorm.DB
	IsRunning bool
}

func Configure(db *gorm.DB) *mux.Router {
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.Common)
	repo := repos.RepositorySQL{Db: db}
	handlers.AddUserRoutes(r, repo)
	handlers.AddSessionRoutes(r, repo)
	handlers.AddTeamRoutes(r, repo)
	return r
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Team{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&models.Player{})
	if err != nil {
		return err
	}
	return nil
}

func CreateApp(address, host, user, password, dbname, port string) (*App, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db")
		return nil, err
	}

	err = Migrate(db)
	if err != nil {
		log.Fatal("Failed to migrate db")
		return nil, err
	}

	app := App{}
	app.address = address
	app.db = db
	app.router = Configure(db)
	return &app, nil
}

func (a *App) Run() {
	http.Handle("/", a.router)
	log.Printf("Listening on address %s\n", a.address)
	l, err := net.Listen("tcp", a.address)
	if err != nil {
		log.Fatal(err)
		return
	}
	a.IsRunning = true
	log.Fatal(http.Serve(l, nil))
}

func (a *App) Close() {
	sqlDB, err := a.db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	_ = sqlDB.Close()
}
