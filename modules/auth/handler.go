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

	if r.RequestURI == config.Vars.Auth.Path.Logout {
		log := ctx.GetLogger("logout")

		log.Debugf("start: %s", config.Vars.Auth.Type)

		switch config.Vars.Auth.Type {
		case "basicauth":
			basicAuthLogout(r, w)
		}

		http.Redirect(w, r, config.Vars.Auth.Path.LogoutRedirect, http.StatusTemporaryRedirect)
		return true

	}

	if r.RequestURI == config.Vars.Auth.Path.Login {
		log := ctx.GetLogger("login")

		log.Debugf("start: %s", config.Vars.Auth.Type)
		var isValid bool
		var handled bool

		switch config.Vars.Auth.Type {
		case "basicauth":
			isValid, handled = basicAuth(r, w)
		}

		if handled {
			return true
		}

		if isValid {

			defaultPath := defaultPath(r)
			resetDefaultPath(w)

			http.Redirect(w, r, defaultPath, http.StatusTemporaryRedirect)
			return true
		}

	}

	return false
}

func basicAuthLogout(r *http.Request, w http.ResponseWriter) {

	r.SetBasicAuth("", "")
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Auth,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)

}
func basicAuth(r *http.Request, w http.ResponseWriter) (bool, bool) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		requireBasicAuth(w)
		return false, true
	}
	isValid := user == config.Vars.Auth.BasicAuth.Username && pass == config.Vars.Auth.BasicAuth.Password
	if !isValid {
		// retry
		requireBasicAuth(w)
		return false, true
	}

	createAuthSession(w)
	return isValid, false
}

func createAuthSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  config.Vars.Auth.Cookies.Auth,
		Value: "admin",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func resetDefaultPath(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Redirect,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, cookie)
}

func defaultPath(r *http.Request) string {
	defaultPath := "/"
	defaultCookie, _ := r.Cookie(config.Vars.Auth.Cookies.Redirect)
	if defaultCookie != nil {
		defaultPath = defaultCookie.Value
	}
	return defaultPath
}

func requireBasicAuth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Unauthorized\n"))
}
