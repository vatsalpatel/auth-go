package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))

		// Extract the JWT token from the Authorization header
		cookies := r.Cookies()
		var tokenString string

		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Remove "Bearer " prefix from the token string
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next middleware or handler
		next.ServeHTTP(w, r)
	})
}
