package routers

import (
	"AuthInGo/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router interface {
	Register(r chi.Router)
}


func SetupRouter(UserRouter Router) *chi.Mux {
    chirouter := chi.NewRouter()

    // Built-in Chi middleware for logging requests
    chirouter.Use(middleware.Logger) 

    chirouter.Get("/ping", controllers.PingHandler)

    UserRouter.Register(chirouter)

    return chirouter
}