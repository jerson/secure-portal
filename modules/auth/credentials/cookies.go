package credentials

import (
	"net/http"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"time"
)

// CookiesCredentials ...
type CookiesCredentials struct {
	Credentials
	Context context.Context

	request *http.Request
	writer  http.ResponseWriter
}

// NewCookiesCredentials ...
func NewCookiesCredentials(context context.Context, r *http.Request, w http.ResponseWriter) *CookiesCredentials {
	return &CookiesCredentials{Context: context, request: r, writer: w}
}

// RedirectPath ...
func (b *CookiesCredentials) RedirectPath() string {
	defaultPath := "/"
	defaultCookie, _ := b.request.Cookie(config.Vars.Auth.Cookies.Redirect)
	if defaultCookie != nil {
		defaultPath = defaultCookie.Value
	}
	return defaultPath
}

// ResetRedirectPath ...
func (b *CookiesCredentials) ResetRedirectPath() {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Redirect,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(b.writer, cookie)
}

// Save ...
func (b *CookiesCredentials) Save() {
	cookie := &http.Cookie{
		Name:  config.Vars.Auth.Cookies.Auth,
		Value: "admin",
		Path:  "/",
	}
	http.SetCookie(b.writer, cookie)
}

// Reset ...
func (b *CookiesCredentials) Reset() {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Auth,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(b.writer, cookie)
}
