package main

import (
	"auth-service/middleware"
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

	mux.Post("/auth/login", app.Authenticate)
	mux.Post("/auth/validate", app.IsTokenBlackListed)
	mux.Post("/auth/refresh", app.RefreshToken)
	mux.Post("/auth/blacklist", app.AddTokenToBlackList)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.AuthMiddleware)

		mux.Post("/all-tokens/{id}", app.OneToken)
	})

	mux.Post("/auth/logout", app.AuthMiddleware(app.Logout))

	return mux
}
