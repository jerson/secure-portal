package auth

import (
	"net/http"
	"secure-portal/modules/config"
)

// IsFirstLoad ...
func IsFirstLoad(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie(config.Vars.Auth.Cookies.Redirect)
	if err != nil {
		cookie := &http.Cookie{
			Name:  config.Vars.Auth.Cookies.Redirect,
			Value: r.RequestURI,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		return true
	}

	return false
}

// IsValid ...
func IsValid(r *http.Request) bool {
	auth, err := r.Cookie(config.Vars.Auth.Cookies.Auth)
	if err != nil {
		return false
	}
	return auth.Value == "admin"
}
