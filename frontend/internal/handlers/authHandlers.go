package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"bitbucket.org/vatsal64/frontend/internal/helpers"
	"bitbucket.org/vatsal64/frontend/internal/models"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	BACKEND_URL := os.Getenv("BACKEND_URL")
	res, err := http.Post(fmt.Sprintf("%v/api/login", BACKEND_URL), "application/json", bytes.NewBuffer([]byte(`{"email": "`+email+`", "password": "`+password+`"}`)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, map[string]interface{}{
			"Error": "Internal server error, please try again later",
		})
		return
	}
	token, ok := data["access_token"]
	if !ok {
		log.Println("No access token found")
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, map[string]interface{}{
			"Error": "Invalid email or password",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: token.(string),
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("BACKEND_URL")+"/api/google/login", http.StatusSeeOther)
}

func ServeSignupPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	BACKEND_URL := os.Getenv("BACKEND_URL")
	res, err := http.Post(fmt.Sprintf("%v/api/signup", BACKEND_URL), "application/json", bytes.NewBuffer([]byte(`{"email": "`+email+`", "password": "`+password+`"}`)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Println("failed to decode response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, ok := data["access_token"]
	if !ok {
		log.Println("No access token found")
		tmpl, _ := template.ParseFiles("templates/signup.html")
		tmpl.Execute(w, map[string]interface{}{
			"Error": "Invalid email or password",
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: token.(string),
	})

	http.Redirect(w, r, "/profile/edit", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "access_token",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	token := helpers.ExtractTokenFromCookies(r.Cookies())
	if token == "" {
		http.Error(w, "No access token found", http.StatusUnauthorized)
		return
	}
	log.Println(token)

	BACKEND_URL := os.Getenv("BACKEND_URL")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/api/profile", BACKEND_URL), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data models.User
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error()+"wew", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func ServeEditProfilePage(w http.ResponseWriter, r *http.Request) {
	token := helpers.ExtractTokenFromCookies(r.Cookies())
	var data models.User
	if token != "" {
		BACKEND_URL := os.Getenv("BACKEND_URL")
		req, _ := http.NewRequest("GET", fmt.Sprintf("%v/api/profile", BACKEND_URL), nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error()+"wew", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/editProfile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	data := models.User{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	body, _ := json.Marshal(data)
	token := helpers.ExtractTokenFromCookies(r.Cookies())

	BACKEND_URL := os.Getenv("BACKEND_URL")
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/api/profile", BACKEND_URL), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewDecoder(res.Body).Decode(&data)

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
