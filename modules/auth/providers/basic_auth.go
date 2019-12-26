package providers

import (
	"fmt"
	"net/http"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"time"
)

// BasicAuthProvider ...
type BasicAuthProvider struct {
	authProvider
	Context context.Context
	Data    map[string]interface{}
	request *http.Request
	writer  http.ResponseWriter
}

// NewBasicAuthProvider ...
func NewBasicAuthProvider(context context.Context, r *http.Request, w http.ResponseWriter) *BasicAuthProvider {
	ctx := &BasicAuthProvider{Context: context, request: r, writer: w}
	ctx.init()
	return ctx
}

// init ...
func (b *BasicAuthProvider) init() {
	b.Data = map[string]interface{}{}
}

// Logout ...
func (b *BasicAuthProvider) Logout() (handled bool) {

	b.request.SetBasicAuth("", "")
	b.resetAuth()

	return false
}

// Login ...
func (b *BasicAuthProvider) Login() (isAuth bool, handled bool) {

	user, pass, ok := b.request.BasicAuth()
	if !ok {
		b.requireBasicAuth()
		return false, true
	}
	isValid := user == config.Vars.Auth.BasicAuth.Username && pass == config.Vars.Auth.BasicAuth.Password
	if !isValid {
		// retry
		b.requireBasicAuth()
		return false, true
	}

	b.createAuthSession()
	return isValid, false
}

func (b *BasicAuthProvider) resetAuth() {
	cookie := &http.Cookie{
		Name:    config.Vars.Auth.Cookies.Auth,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}
	http.SetCookie(b.writer, cookie)
}

func (b *BasicAuthProvider) requireBasicAuth() {
	b.writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
	b.writer.WriteHeader(http.StatusUnauthorized)
	b.writer.Write([]byte("401 Unauthorized\n"))
}

func (b *BasicAuthProvider) createAuthSession() {
	cookie := &http.Cookie{
		Name:  config.Vars.Auth.Cookies.Auth,
		Value: "admin",
		Path:  "/",
	}
	http.SetCookie(b.writer, cookie)
}
