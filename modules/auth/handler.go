package auth

import (
	"fmt"
	"net/http"
	"secure-portal/modules/auth/credentials"
	"secure-portal/modules/auth/providers"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
)

// Handler ...
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

	isValidRequest := r.Method == http.MethodGet
	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Logout {
		log := ctx.GetLogger("logout")

		log.Debugf("start: %s", config.Vars.Auth.Type)

		var provider providers.AuthProvider
		credential := credentials.NewCookiesCredentials(ctx, r, w)

		switch config.Vars.Auth.Type {
		case "basicauth":
			provider = providers.NewBasicAuthProvider(ctx, credential, r, w)
		default:
			panic(fmt.Sprintf("type not handled: %s", config.Vars.Auth.Type))
		}

		handled := provider.Logout()
		if handled {
			return true
		}

		credential.ResetRedirectPath()

		http.Redirect(w, r, config.Vars.Auth.Path.LogoutRedirect, http.StatusTemporaryRedirect)
		return true

	}

	if isValidRequest && r.RequestURI == config.Vars.Auth.Path.Login {
		log := ctx.GetLogger("login")

		log.Debugf("start: %s", config.Vars.Auth.Type)

		var provider providers.AuthProvider
		credential := credentials.NewCookiesCredentials(ctx, r, w)

		switch config.Vars.Auth.Type {
		case "basicauth":
			provider = providers.NewBasicAuthProvider(ctx, credential, r, w)
		default:
			panic(fmt.Sprintf("type not handled: %s", config.Vars.Auth.Type))
		}

		isValid, handled := provider.Login()
		if handled {
			return true
		}

		if isValid {
			redirectPath := credential.RedirectPath()
			credential.ResetRedirectPath()

			http.Redirect(w, r, redirectPath, http.StatusTemporaryRedirect)
			return true
		}

	}

	return false
}
