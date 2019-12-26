package auth

import (
	"net/http"
)

func IsFirstLoad(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("Default")
	if err != nil {
		cookie := &http.Cookie{
			Name:  "Default",
			Value: r.RequestURI,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		return true
	}

	return false
}

func IsValid(r *http.Request) bool {
	auth, err := r.Cookie("Auth-Portal")
	if err != nil {
		return false
	}
	return auth.Value == "admin"
}