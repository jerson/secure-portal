package auth

import (
	"net/http"
	"secure-portal/modules/auth/providers"
	"secure-portal/modules/config"
)

// Handler ...
func Handler(provider providers.AuthProvider, w http.ResponseWriter, r *http.Request) (handled bool) {

	isValidRequest := r.Method == http.MethodGet || r.Method == http.MethodPost

	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Logout {

		handled := provider.Logout()
		if handled {
			return true
		}

		provider.Session().ResetRedirectPath()
		http.Redirect(w, r, config.Vars.Auth.Path.LogoutRedirect, http.StatusTemporaryRedirect)
		return true

	}

	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Register {

		handled := provider.Register()
		if handled {
			return true
		}

		http.Redirect(w, r, config.Vars.Auth.Path.Login, http.StatusTemporaryRedirect)
		return true

	}

	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Login {

		isValid, handled := provider.Login()
		if handled {
			return true
		}

		if isValid {
			redirectPath := provider.Session().RedirectPath()
			if redirectPath == "" {
				redirectPath = "/"
			}
			provider.Session().ResetRedirectPath()

			http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
			return true
		}

	}

	return false
}
