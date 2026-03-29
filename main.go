package main

import (
	"challenge2/internal/handlers"
	"challenge2/internal/middleware"
	"challenge2/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	authService := service.NewAuthService()
	authHandler := handlers.NewAuthHandler(authService)

	r.Post("/api/user/register", authHandler.RegisterHandler)
	r.Post("/api/user/login", authHandler.LoginHandler)
	r.With(middleware.AuthMiddleware(authService)).Get("/api/user/profile", authHandler.ProfileHandler)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
