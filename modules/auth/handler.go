package auth

import (
	"net/http"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request){
	if r.RequestURI == "/auth" {
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		isValidCredentials := user == "admin" && pass == "admin"
		if !isValidCredentials {
			// retry
			w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		if isValidCredentials {

			cookie := &http.Cookie{
				Name:  "Auth-Portal",
				Value: "admin",
				Path:  "/",
			}
			http.SetCookie(w, cookie)

			defaultPath := "/"
			defaultCookie, _ := r.Cookie("Default")
			if defaultCookie != nil {
				defaultPath = defaultCookie.Value
			}

			resetDefaultCookie := &http.Cookie{
				Name:    "Default",
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),
			}
			http.SetCookie(w, resetDefaultCookie)

			http.Redirect(w, r, defaultPath, http.StatusTemporaryRedirect)
			return
		}
	}
}