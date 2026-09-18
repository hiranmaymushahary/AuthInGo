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

// SetupRouter initializes the router, configures middleware, and registers module routes
func SetupRouter(UserRouter Router) *chi.Mux {
    // Standardized casing to chiRouter (camelCase)
    chiRouter := chi.NewRouter()

    // Built-in Chi middleware for logging requests
    chiRouter.Use(middleware.Logger)   

    chiRouter.Use(middlewares.RateLimitMiddleware)
   

    // Route for basic server health check
    chiRouter.Get("/ping", controllers.PingHandler)

    // Register user routes onto the main chi router
    UserRouter.Register(chiRouter)

    return chiRouter
}