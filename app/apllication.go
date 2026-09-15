package app

import (
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
}

// CONSTRUCTER FRO CONFIG

func NewConfig(addr string) Config {
	return Config{
		Addr: addr,
	}

}

// CONSTRUCTER FOR APPLICATION

func NewApplication(cfg Config) *Application {
	return &Application{
		Config: cfg,
	}
}

func (app *Application) Run() error {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Strating server on", app.Config.Addr)

	return server.ListenAndServe()
}
