package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/vatsal64/va_pa/internal/handlers"
	"bitbucket.org/vatsal64/va_pa/internal/helpers"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware_ValidToken(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-jwt-middleware@example.com")

	// Create a user in the database.
	user, err := helpers.RegisterUser("test-jwt-middleware@example.com", "password123", "profile", "123456", false)
	if err != nil {
		t.Fatal(err)
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request with JSON data for a signup with a duplicate email.
	req, err := http.NewRequest("GET", "/profile", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Call the ProfileHandler to perform the signup.
	JWTMiddleware(mockHandler).ServeHTTP(rr, req)

	// Call the ProfileHandler to perform the signup.
	handlers.ProfileHandler(rr, req)

	db.Exec("DELETE FROM users WHERE email = ?", "test-jwt-middleware@example.com")

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

}

func TestJWTMiddleware_InvalidToken(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-jwt-middleware@example.com")

	// Create a user in the database.
	user, err := helpers.RegisterUser("test-jwt-middleware@example.com", "password123", "profile", "123456", false)
	if err != nil {
		t.Fatal(err)
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}
	token = token + "invalid"

	// Create a request with JSON data for a signup with a duplicate email.
	req, err := http.NewRequest("GET", "/profile", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Call the ProfileHandler to perform the signup.
	JWTMiddleware(mockHandler).ServeHTTP(rr, req)

	db.Exec("DELETE FROM users WHERE email = ?", "test-jwt-middleware@example.com")

	// Check the status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
