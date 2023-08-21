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

	email, ok := reqData["email"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "invalid email")))
		return
	}
	password, ok := reqData["password"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "invalid password")))
		return
	}

	user, err := helpers.ValidateLogin(email, password)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "email or password is incorrect")))
		return
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "email or password is incorrect")))
		return
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
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

	email, ok := reqData["email"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "invalid email")))
		return
	}
	password, ok := reqData["password"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "invalid password")))
		return
	}

	user, err := helpers.RegisterUser(email, password, "", "", false)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, "email is already registered")))
		return
	}

	token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	user := helpers.GetUserByID(claims.UserID)
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	w.Write(body)
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var reqData models.User
	json.NewDecoder(r.Body).Decode(&reqData)

	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	claims, err := helpers.ExtractClaims(token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	reqData.ID = claims.UserID
	err = helpers.UpdateUserByID(reqData.ID, reqData.Name, reqData.Email, reqData.Phone)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	user := helpers.GetUserByID(claims.UserID)
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	w.Write(body)
}

func GetGoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := config.GoogleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusSeeOther)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	email := profile["email"].(string)
	user, err := helpers.ValidateGoogleLogin(email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}
	// User.ID will be 0 if user has not logged in with google before, then register them with email as username
	if user.ID == 0 {
		user, err = helpers.RegisterUser(email, "", profile["name"].(string), "", true)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			return
		}
	}
	access_token, err := helpers.GenerateToken(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		return
	}

	// Normally there would be a redirect response with a set cookie header, but due to domain being in public suffix list
	// I have opted for performing the exchange using javascript in another window
	// So access_token is provided as response which will be passed to frontend which is the opener of this window, and this exchange window will be closed
	w.Write([]byte(fmt.Sprintf(`<script>window.opener.postMessage('{"access_token": "%v"}', '*'); window.close();</script>`, access_token)))

	// This is what it would look like normally
	// http.SetCookie(w, &http.Cookie{
	// 	Name:  "access_token",
	// 	Value: access_token,
	// 	Path:  "/",
	// 	Domain: os.Getenv("COOKIE_DOMAIN"),
	// })
	// FRONTEND_URI := os.Getenv("FRONTEND_URI")
	// http.Redirect(w, r, FRONTEND_URI+"/login/google/success", http.StatusFound)
}
