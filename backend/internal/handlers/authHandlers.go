package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"bitbucket.org/vatsal64/va_pa/config"
	"bitbucket.org/vatsal64/va_pa/internal/helpers"
	"bitbucket.org/vatsal64/va_pa/internal/models"
	"golang.org/x/oauth2"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	var reqData map[string]string
	json.Unmarshal(body, &reqData)

	username, ok := reqData["username"]
	if !ok {
		w.Write([]byte(`{"error": "invalid username"}`))
		return
	}
	password, ok := reqData["password"]
	if !ok {
		w.Write([]byte(`{"error": "invalid password"}`))
		return
	}

	user, err := helpers.ValidateLogin(username, password)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}

	if user.ID == 0 {
		w.Write([]byte(`{"error": "username or password is incorrect`))
		return
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: token,
	})

	resp, _ := json.Marshal(map[string]interface{}{
		"message":      "login successful",
		"access_token": token,
	})
	w.Write(resp)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)

	var reqData map[string]string
	json.Unmarshal(body, &reqData)

	username, ok := reqData["username"]
	if !ok {
		w.Write([]byte(`{"error": "invalid username"}`))
		return
	}
	password, ok := reqData["password"]
	if !ok {
		w.Write([]byte(`{"error": "invalid password"}`))
		return
	}

	user, err := helpers.RegisterUser(username, password, "", "", "")
	if err != nil {
		log.Println(err)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "Username is already registered")))
		return
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: token,
	})

	resp, _ := json.Marshal(map[string]interface{}{
		"message":      "signup successful",
		"access_token": token,
	})
	w.Write(resp)
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := helpers.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := helpers.GetUserByID(claims.UserID)
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(body))
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var reqData models.User
	json.NewDecoder(r.Body).Decode(&reqData)
	log.Println(reqData)

	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := helpers.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqData.ID = claims.UserID
	err = helpers.UpdateUserByID(reqData.ID, reqData.Name, reqData.Email, reqData.Phone)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := helpers.GetUserByID(claims.UserID)
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func GetGoogleLoginURLHandler(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	w.Write([]byte(url))
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := helpers.ValidateLogin(email, "")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.ID == 0 {
		user, err = helpers.RegisterUser(email, "", profile["name"].(string), email, "") // Username is same as email for google login
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	access_token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: access_token,
	})

	responseBody, err := json.Marshal(map[string]interface{}{
		"message":      "login successful",
		"access_token": access_token,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(responseBody))
}
