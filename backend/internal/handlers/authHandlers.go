package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"bitbucket.org/vatsal64/va_pa/config"
	"golang.org/x/oauth2"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("signup"))
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("profile"))
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("edit profile"))
}

func GetGoogleLoginURLHandler(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	w.Write([]byte(url))
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("google callback")
	code := r.URL.Query().Get("code")
	token, err := config.GoogleOauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := config.GoogleOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := profile["email"].(string)
	log.Println(profile)

	http.SetCookie(w, &http.Cookie{
		Name:  "email",
		Value: email,
	})
	responseBody, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseBody))
}
