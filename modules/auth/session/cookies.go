package session

import (
	"net/http"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"time"
)

// CookiesSession ...
type CookiesSession struct {
	Session
	Context context.Context

	request *http.Request
	writer  http.ResponseWriter
}

// NewCookiesSession ...
func NewCookiesSession(context context.Context, r *http.Request, w http.ResponseWriter) *CookiesSession {
	return &CookiesSession{Context: context, request: r, writer: w}
}

// RedirectPath ...
func (s *CookiesSession) RedirectPath() string {
	defaultPath := ""
	defaultCookie, _ := s.request.Cookie(config.Vars.Auth.Cookies.Redirect)
	if defaultCookie != nil {
		defaultPath = defaultCookie.Value
	}
	return defaultPath
}

// SetRedirectPath ...
func (s *CookiesSession) SetRedirectPath(path string) {
	cookie := &http.Cookie{
		Name:  config.Vars.Auth.Cookies.Redirect,
		Value: path,
		Path:  "/",
	}
	http.SetCookie(s.writer, cookie)
}

// GetToken ...
func (s *CookiesSession) GetToken() string {
	auth, err := s.request.Cookie(config.Vars.Auth.Cookies.Auth)
	if err != nil {
		return ""
	}
	return auth.Value
}

// ResetRedirectPath ...
func (s *CookiesSession) ResetRedirectPath() {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Redirect,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(s.writer, cookie)
}

// Save ...
func (s *CookiesSession) Save(token string) {
	cookie := &http.Cookie{
		Name:  config.Vars.Auth.Cookies.Auth,
		Value: token,
		Path:  "/",
	}
	http.SetCookie(s.writer, cookie)
}

// Reset ...
func (s *CookiesSession) Reset() {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Auth,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(s.writer, cookie)
}
