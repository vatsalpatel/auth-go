package routes

import (
	"bitbucket.org/vatsal64/frontend/internal/handlers"
	"bitbucket.org/vatsal64/frontend/internal/middleware"
	"github.com/go-chi/chi"
)

func Configure(r *chi.Mux) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RedirectMiddleware)

		r.Get("/login", handlers.ServeLoginPage)
		r.Post("/login", handlers.LoginHandler)
		r.Get("/login/google", handlers.GoogleLoginHandler)

		r.Get("/signup", handlers.ServeSignupPage)
		r.Post("/signup", handlers.SignupHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Get("/logout", handlers.LogoutHandler)

		r.Get("/profile", handlers.ProfileHandler)

		r.Get("/profile/edit", handlers.ServeEditProfilePage)
		r.Post("/profile/edit", handlers.EditProfileHandler)
	})
}
