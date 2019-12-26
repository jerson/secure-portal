package credentials

// Credentials ...
type Credentials interface {
	RedirectPath() string
	ResetRedirectPath()
	Save()
	Reset()
}
