package app

import (
	dbConfig "AuthInGo/config/db"
	config "AuthInGo/config/env"
	"AuthInGo/controllers"
	repo "AuthInGo/db/repositories"
	"AuthInGo/routers"
	"AuthInGo/services"
	"fmt"
	"net/http"
	"time"
)

// CONFIG HOLDS THE CONFIGURATION FOR THE SERVER.(SERVER CONFIGURATION)
type Config struct {
	Addr string // Port
}

// CONTAIN SERVER DETAILS (GET ALL THE CONFIGURATION FROM THE CONFIG STRUCT)
type Application struct {
	Config Config
	Store  repo.Storage
}

// CONSTRUCTER FRO CONFIG

func NewConfig() Config {

	port := config.GetString("PORT", ":8080")
	return Config{
		Addr: port,
	}

}

// CONSTRUCTER FOR APPLICATION

func NewApplication(cfg Config) *Application {
	return &Application{
		Config: cfg,
		Store:  *repo.NewStorage(),
	}
}

func (app *Application) Run() error {

	db, err := dbConfig.SetupDB()

	if err != nil {
		fmt.Println("Erro setting up in database")
		return err
	}

	ur := repo.NewUserRepository(db)
	us := services.NewUserService(ur)
	uc := controllers.NewUserController(us)
	uRouter := routers.NewUserRouter(uc)

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      routers.SetupRouter(uRouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Strating server on", app.Config.Addr)

	return server.ListenAndServe()
}
