package providers

import (
	"fmt"
	"net/http"
	"secure-portal/modules/auth/credentials"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
)

// BasicAuthProvider ...
type BasicAuthProvider struct {
	AuthProvider
	Context     context.Context
	Credentials credentials.Credentials
	Data        map[string]interface{}

	request *http.Request
	writer  http.ResponseWriter
}

// NewBasicAuthProvider ...
func NewBasicAuthProvider(context context.Context, credentials credentials.Credentials, r *http.Request, w http.ResponseWriter) *BasicAuthProvider {
	ctx := &BasicAuthProvider{Context: context, Credentials: credentials, request: r, writer: w}
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
	b.Credentials.Reset()

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

	b.Credentials.Save()
	return isValid, false
}

func (b *BasicAuthProvider) requireBasicAuth() {
	b.writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
	b.writer.WriteHeader(http.StatusUnauthorized)
	b.writer.Write([]byte("401 Unauthorized\n"))
}
