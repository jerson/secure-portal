package auth

import (
	"fmt"
	"net/http"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"time"
)

// Handler ...
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

	if r.RequestURI == "/auth" {
		log := ctx.GetLogger("auth")

		log.Debug("start auth")
		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return true
		}
		isValidCredentials := user == "admin" && pass == "admin"
		if !isValidCredentials {
			// retry
			w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized\n"))
			return true
		}
		if isValidCredentials {

			cookie := &http.Cookie{
				Name:  config.Vars.Auth.Cookies.Auth,
				Value: "admin",
				Path:  "/",
			}
			http.SetCookie(w, cookie)

			defaultPath := "/"
			defaultCookie, _ := r.Cookie(config.Vars.Auth.Cookies.Redirect)
			if defaultCookie != nil {
				defaultPath = defaultCookie.Value
			}

			resetDefaultCookie := &http.Cookie{
				Name:    config.Vars.Auth.Cookies.Redirect,
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),
			}
			http.SetCookie(w, resetDefaultCookie)

			http.Redirect(w, r, defaultPath, http.StatusTemporaryRedirect)
			return true
		}
	}

	return false
}
