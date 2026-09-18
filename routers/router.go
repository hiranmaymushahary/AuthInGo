package routers

import (
	"AuthInGo/controllers"
    "AuthInGo/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	Register(r chi.Router)
}


func SetupRouter(UserRouter Router) *chi.Mux {
    chirouter := chi.NewRouter()

    // Built-in Chi middleware for logging requests
    chiRouter.Use(middleware.Logger)   
    // Middleware for validating requests         
	chiRouter.Use(middlewares.RequestValidator) 

    chirouter.Get("/ping", controllers.PingHandler)

    UserRouter.Register(chirouter)

    return chirouter
}