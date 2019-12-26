package providers

// AuthProvider ...
type AuthProvider interface {
	Login() (isAuth bool, handled bool)
	Logout() (handled bool)
}
