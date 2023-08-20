package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/vatsal64/va_pa/internal/helpers"
	"bitbucket.org/vatsal64/va_pa/pkg/storage"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}
	// Create a user in the database.
	_, err := helpers.RegisterUser("test-login@example.com", "password123", "", "", false)
	if err != nil {
		t.Fatal(err)
	}
	// Create a request with JSON data (adjust as needed for your specific JSON structure).
	reqBody := []byte(`{"email": "test-login@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the LoginHandler.
	http.HandlerFunc(LoginHandler).ServeHTTP(rr, req)

	// Check the status code and response body.
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "login successful", responseBody["message"])
	assert.NotEmpty(t, responseBody["access_token"])
}

func TestLoginHandlerIncorrectEmail(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}
	// Create a request with missing email field.
	reqBody := []byte(`{"email": "test@example2.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the LoginHandler.
	http.HandlerFunc(LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "email or password is incorrect", responseBody["error"])
}

func TestLoginHandlerIncorrectPassword(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}
	// Create a request with missing email field.
	reqBody := []byte(`{"email": "test@example.com", "password": "password1234"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the LoginHandler.
	http.HandlerFunc(LoginHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "email or password is incorrect", responseBody["error"])
}

func TestSignupHandler_ValidSignup(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-signup@example.com")

	// Create a request with JSON data for a valid signup.
	reqBody := []byte(`{"email": "test-signup@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the SignupHandler to perform the signup.
	SignupHandler(rr, req)

	// Check the status code and response body.
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	db.Exec("DELETE FROM users WHERE email = ?", "test-signup@example.com")

	assert.Equal(t, "signup successful", responseBody["message"])
	assert.NotNil(t, responseBody["access_token"])
}

func TestSignupHandler_DuplicateEmail(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-duplicate@example.com")

	// Create a user with a duplicate email in the database.
	_, err := helpers.RegisterUser("test-duplicate@example.com", "password123", "", "", false)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request with JSON data for a signup with a duplicate email.
	reqBody := []byte(`{"email": "test-duplicate@example.com", "password": "password123"}`)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the SignupHandler to perform the signup.
	SignupHandler(rr, req)

	// Check the status code and response body.
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	db.Exec("DELETE FROM users WHERE email = ?", "test-duplicate@example.com")

	assert.Equal(t, "email is already registered", responseBody["error"])
}

func TestProfileHandler(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-profile@example.com")

	// Create a user in the database.
	user, err := helpers.RegisterUser("test-profile@example.com", "password123", "profile", "123456", false)
	if err != nil {
		t.Fatal(err)
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request with JSON data for a signup with a duplicate email.
	reqBody := []byte(`{"email": "test-profile@example.com", "password": "password123"}`)
	req, err := http.NewRequest("GET", "/profile", bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the ProfileHandler to perform the signup.
	ProfileHandler(rr, req)

	// Check the status code and response body.
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	db.Exec("DELETE FROM users WHERE email = ?", "test-profile@example.com")

	assert.Equal(t, "test-profile@example.com", responseBody["email"])
	assert.Equal(t, "profile", responseBody["name"])
	assert.Equal(t, "123456", responseBody["phone"])

}

func TestEditProfileHandler(t *testing.T) {
	if godotenv.Load("../../.env.test") != nil {
		t.Fatal("Error loading .env file")
	}

	db := storage.GetMySQLDatabase()
	db.Exec("DELETE FROM users WHERE email = ?", "test-edit-profile@example.com")
	db.Exec("DELETE FROM users WHERE email = ?", "test-edited-profile@example.com")

	// Create a user in the database.
	user, err := helpers.RegisterUser("test-edit-profile@example.com", "password123", "profile", "123456", false)
	if err != nil {
		t.Fatal(err)
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request with JSON data for a signup with a duplicate email.
	reqBody := []byte(`{"email": "test-edited-profile@example.com", "name": "new name", "phone": "654321"}`)
	req, err := http.NewRequest("PUT", "/profile", bytes.NewBuffer(reqBody))
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response.
	rr := httptest.NewRecorder()

	// Call the ProfileHandler to perform the signup.
	EditProfileHandler(rr, req)

	// Check the status code and response body.
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	db.Exec("DELETE FROM users WHERE email = ?", "test-edit-profile@example.com")
	db.Exec("DELETE FROM users WHERE email = ?", "test-edited-profile@example.com")

	assert.Equal(t, "test-edited-profile@example.com", responseBody["email"])
	assert.Equal(t, "new name", responseBody["name"])
	assert.Equal(t, "654321", responseBody["phone"])
}
