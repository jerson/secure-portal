package auth

import (
	"net/http"
	"secure-portal/modules/auth/providers"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"time"
)

// Handler ...
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

	isValidRequest := r.Method == http.MethodGet
	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Logout {
		log := ctx.GetLogger("logout")

		log.Debugf("start: %s", config.Vars.Auth.Type)
		var handled bool

		switch config.Vars.Auth.Type {
		case "basicauth":
			provider := providers.NewBasicAuthProvider(ctx, r, w)
			handled = provider.Logout()
		}

		if handled {
			return true
		}

		resetDefaultPath(w)
		http.Redirect(w, r, config.Vars.Auth.Path.LogoutRedirect, http.StatusTemporaryRedirect)
		return true

	}

	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Login {
		log := ctx.GetLogger("login")

		log.Debugf("start: %s", config.Vars.Auth.Type)
		var isValid bool
		var handled bool

		switch config.Vars.Auth.Type {
		case "basicauth":
			provider := providers.NewBasicAuthProvider(ctx, r, w)
			isValid, handled = provider.Login()
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
