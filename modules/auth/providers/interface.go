package providers

import "secure-portal/modules/auth/session"

// AuthProvider ...
type AuthProvider interface {
	Login() (isAuth bool, handled bool)
	Logout() (handled bool)
	Session() (session session.Session)
	IsAuthenticated() (success bool)
	IsFirstTime() (success bool)
}
