package providers

import (
	"net/http"
)

// authProvider ...
type authProvider interface {
	Login(r *http.Request, w http.ResponseWriter) (isAuth bool, handled bool)
	Logout(r *http.Request, w http.ResponseWriter) (handled bool)
}
