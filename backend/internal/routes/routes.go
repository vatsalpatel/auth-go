package routes

import (
	"bitbucket.org/vatsal64/va_pa/internal/handlers"
	"bitbucket.org/vatsal64/va_pa/internal/middleware"
	"github.com/go-chi/chi"
)

func Configure(r *chi.Mux) {
	r.Post("/api/login", handlers.LoginHandler)
	r.Post("/api/signup", handlers.SignupHandler)
	r.Get("/api/google/login", handlers.GetGoogleLoginHandler)
	r.Get("/api/callback", handlers.GoogleCallbackHandler)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Get("/api/profile", handlers.ProfileHandler)
		r.Put("/api/profile", handlers.EditProfileHandler)
	})
}
