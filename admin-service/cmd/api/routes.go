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

	userRouter := chi.NewRouter()
	mux.Mount("/user", userRouter)

	userRouter.Post("/register", app.Register)
	userRouter.Get("/check-email", app.CheckEmailExists)
	userRouter.Post("/forgot-password", app.SendPasswordResetEmail)
	userRouter.Post("/reset-password", app.ResetPassword)

	userRouter.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) 

		r.Delete("/delete/{user_id}", app.DeleteUserByID)
		r.Post("/change-password", app.ChangePassword)
		r.Put("/update/{user_id}", app.UpdateUser)
	})

	return mux
}
