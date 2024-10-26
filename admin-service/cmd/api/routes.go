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

	mux.Post("/api/admin/register", app.Register)
	mux.Post("api/admin/reset-password", app.ResetPassword)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.AuthMiddleware("admin"))

		mux.Post("/new-add", app.AddNewAdmin)
		mux.Get("/new-admins", app.GetAllNewAdmins)
		mux.Delete("/delete/{admin_id}", app.DeleteAdmin)
		mux.Delete("/delete-new/", app.DeleteNewAdmin)
		mux.Put("/update/{admin_id}", app.UpdateAdmin)

	})

	return mux
}
