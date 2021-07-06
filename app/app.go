package app

import (
	"./controllers"
	"./middleware"
	"./migrations"
	"./repos"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net"
	"net/http"
)

type App struct {
	address   string
	router    *mux.Router
	db        *gorm.DB
	IsRunning bool
	user *controllers.UserController
	session *controllers.SessionController
	player *controllers.PlayerController
	team *controllers.TeamController
	transfer *controllers.TransferController
}

func (a *App) Configure() {
	repo := repos.RepositorySQL{Db: a.db}
	a.user = &controllers.UserController{Repo: repo}
	a.session = &controllers.SessionController{Repo: repo}
	a.player = &controllers.PlayerController{Repo: repo}
	a.team = &controllers.TeamController{Repo: repo}
	a.transfer = &controllers.TransferController{Repo: repo}

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.Common)

	a.user.AddRoutes(r)
	a.session.AddRoutes(r)
	a.team.AddRoutes(r)
	a.player.AddRoutes(r)
	a.transfer.AddRoutes(r)

	a.router = r
}

func CreateApp(address, host, user, password, dbname, port string) (*App, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("Failed to connect to db")
		return nil, err
	}

	err = migrations.Run(db)
	if err != nil {
		log.Fatal("Failed to migrate db")
		return nil, err
	}

	app := App{}
	app.address = address
	app.db = db
	app.Configure()
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
