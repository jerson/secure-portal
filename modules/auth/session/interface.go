package session

// Session ...
type Session interface {
	RedirectPath() string
	SetRedirectPath(path string)
	ResetRedirectPath()
	Save(token string)
	GetToken() string
	Reset()
}
