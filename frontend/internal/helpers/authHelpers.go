package helpers

import "net/http"

func ExtractTokenFromCookies(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == "access_token" {
			return cookie.Value
		}
	}
	return ""
}

func IsLoggedIn(cookies []*http.Cookie) bool {
	return ExtractTokenFromCookies(cookies) != ""
}
