package middleware

import (
	"net/http"

	"bitbucket.org/vatsal64/frontend/internal/helpers"
)

func RedirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := helpers.ExtractTokenFromCookies(r.Cookies())
		if token != "" {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
