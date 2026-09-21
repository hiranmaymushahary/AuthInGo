package routers

import (
    "AuthInGo/controllers"
    "AuthInGo/middlewares"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "AuthInGo/utils"
)

type Router interface {
    Register(r chi.Router)
}

// SetupRouter initializes the router, configures middleware, and registers module routes
func SetupRouter(UserRouter Router , RoleRouter Router) *chi.Mux {
    // Standardized casing to chiRouter (camelCase)
    chiRouter := chi.NewRouter()

    // Built-in Chi middleware for logging requests
    chiRouter.Use(middleware.Logger)   

    chiRouter.Use(middlewares.RateLimitMiddleware)
   

    // Route for basic server health check
    chiRouter.Get("/ping", controllers.PingHandler)


    chiRouter.HandleFunc("/fakestoreservice/*", utils.ProxyToService("https://fakestoreapi.in", "/fakestoreservice"))


    // Register user routes onto the main chi router
    UserRouter.Register(chiRouter)

    RoleRouter.Register(chiRouter)

    return chiRouter
}

// http://localhost:3001/fakestoreservice/products/category