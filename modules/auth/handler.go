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

		log.Debugf("start auth: %s", config.Vars.Auth.Type)
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

func basicAuth(r *http.Request, w http.ResponseWriter) (bool, bool) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		requireBasicAuth(w)
		return false, true
	}
	isValid := user == "admin" && pass == "admin"
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
	resetDefaultCookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Redirect,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, resetDefaultCookie)
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
