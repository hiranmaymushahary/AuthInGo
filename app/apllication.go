package app

import (
	config "AuthInGo/config/env"
	"AuthInGo/controllers"
	db "AuthInGo/db/repositories"
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
	Store  db.Storage
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
		Store:  *db.NewStorage(),
	}
}

func (app *Application) Run() error {

	ur := db.NewUserRepository()
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
