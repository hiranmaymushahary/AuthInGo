package routers

import (
	"AuthInGo/controllers"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	Register(r chi.Router)
}

func SetupRouter(UserRouter Router) *chi.Mux {
	chirouter := chi.NewRouter()

	chirouter.Get("/ping", controllers.PingHandler)

	UserRouter.Register(chirouter)

	return chirouter
}
