package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	// Create a request with JSON data (adjust as needed for your specific JSON structure).
	reqBody := []byte(`{"email": "test@example.com", "password": "password123"}`)
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
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"message": "login successful", "access_token": "your_access_token_here"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestLoginHandlerInvalidEmail(t *testing.T) {
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

	// Check the status code and response body.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"error": "invalid email"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestLoginHandlerInvalidPassword(t *testing.T) {
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

	// Check the status code and response body.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"error": "invalid email"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}
