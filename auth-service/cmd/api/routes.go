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

	mux.Post("/api/auth/login", app.AuthenticateUser)
	mux.Post("/api/admin/login", app.AuthenticateAdmin)
	mux.Post("/api/auth/refresh", app.RefreshToken) 

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.AuthMiddleware("admin"))
		
		mux.Post("/api/auth/revoke", app.RevokeToken)
		//mux.Post("/all-tokens/{id}", app.OneToken) 
	})

	mux.Post("/auth/logout", app.AuthMiddleware("user")(http.HandlerFunc(app.Logout)).ServeHTTP)

	return mux
}
