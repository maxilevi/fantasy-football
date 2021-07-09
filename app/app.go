package app

import (
	_ "../docs"
	"./controller"
	"./middleware"
	"./migrations"
	"./repos"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net"
	"net/http"
)

type App struct {
	address   string
	db        *gorm.DB
	router    *gin.Engine
	IsRunning bool
}

func (a *App) Configure() {
	repo := repos.RepositorySQL{Db: a.db}
	r := gin.Default()

	c := controller.NewController(repo)

	api := r.Group("/api")
	{
		users := api.Group("/user")
		{
			users.GET("/:id", c.ShowUser)
			users.POST("", c.CreateUser)
			users.Use(middleware.Auth(repo))
			users.GET("/me", c.ShowMyself)
			users.Use(middleware.Admin())
			users.DELETE("/:id", c.DeleteUser)
			users.PATCH("/:id", c.UpdateUser)
		}
		session := api.Group("/session")
		{
			session.POST("", c.CreateSession)
		}
		team := api.Group("/team")
		{
			team.GET("/:id", c.ShowTeam)
			team.Use(middleware.Auth(repo))
			team.PATCH("/:id", c.UpdateTeam)
			team.Use(middleware.Admin())
			team.POST("", c.CreateTeam)
			team.DELETE("/:id", c.DeleteTeam)
		}
		players := api.Group("/player")
		{
			players.GET("/:id", c.ShowPlayer)
			players.Use(middleware.Auth(repo))
			players.PATCH("/:id", c.UpdatePlayer)
			players.Use(middleware.Admin())
			players.POST("", c.CreatePlayer)
			players.DELETE("/:id", c.DeletePlayer)
		}
	}
	url := ginSwagger.URL("http://" + a.address + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
