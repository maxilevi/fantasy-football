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

// A struct holding of our server info
type App struct {
	address   string
	db        *gorm.DB
	router    *gin.Engine
	IsRunning bool
}

// Configure the app routes and its data source
func (a *App) Configure() {
	repo := repos.RepositorySQL{Db: a.db}
	r := gin.Default()

	c := controller.NewController(repo)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("", c.CreateUser)
			users.Use(middleware.Auth(repo))
			users.GET("/me", c.RedirectMyself)
			users.GET("/me/team", c.GetMyTeam)
			users.PATCH("/me/team", c.EditMyTeam)
			users.GET("/me/team/players", c.GetMyPlayers)
			users.GET("/me/team/players/:playerId", c.GetMyPlayer)
			users.PATCH("/me/team/players/:playerId", c.EditMyPlayer)
			users.GET("/:userId", c.ShowUser)
			users.Use(middleware.Admin())
			users.DELETE("/:userId", c.DeleteUser)
			users.PATCH("/:userId", c.UpdateUser)
		}
		session := api.Group("/sessions")
		{
			session.POST("", c.CreateSession)
		}
		team := api.Group("/teams")
		{
			team.GET("/:teamId/players", c.ListTeamPlayers)
			team.GET("/:teamId/players/:playerId", c.GetMyPlayerFromTeam)
			team.PATCH("/:teamId/players/:playerId", c.EditMyPlayerFromTeam)
			team.GET("/:teamId", c.ShowTeam)
			team.Use(middleware.Auth(repo))
			team.PATCH("/:teamId", c.UpdateTeam)
			team.Use(middleware.Admin())
			team.POST("/:teamId/players", c.CreateNewPlayerOnTeam)
			team.POST("", c.CreateTeam)
			team.DELETE("/:teamId", c.DeleteTeam)
		}
		players := api.Group("/players")
		{
			players.GET("/:playerId", c.ShowPlayer)
			players.Use(middleware.Auth(repo))
			players.PATCH("/:playerId", c.UpdatePlayer)
			players.Use(middleware.Admin())
			players.DELETE("/:playerId", c.DeletePlayer)
		}
		transfers := api.Group("/transfers")
		{
			transfers.GET("", c.ListTransfers)
			transfers.GET("/:transferId", c.ShowTransfer)
			transfers.Use(middleware.Auth(repo))
			transfers.DELETE("/:transferId", c.DeleteTransfer)
			transfers.PATCH("/:transferId", c.UpdateTransfer)
			transfers.POST("", c.CreateTransfer)
			transfers.PUT("/:transferId/buy", c.BuyTransfer)
		}
	}
	url := ginSwagger.URL("http://" + a.address + "/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	a.router = r
}

// Create a new app with the given parameters
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

// Listen on the specified address and serve requests
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

// Close the app and all it's resources
func (a *App) Close() {
	sqlDB, err := a.db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	_ = sqlDB.Close()
}
