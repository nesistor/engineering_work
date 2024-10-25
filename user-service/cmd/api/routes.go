package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "PUT", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/api/login/register", app.Register)
	mux.Get("/api/login/check-email", app.CheckEmail)
	mux.Post("api/login/reset-password", app.ResetPassword)

	mux.Route("/api/login", func(mux chi.Router) {
		mux.Use(app.AuthMiddleware("user"))
		
		mux.Delete("/delete-user/{user_id}", app.DeleteUser)
		mux.Put("/update/{user_id}", app.UpdateUser)
		
	})

	return mux
}
