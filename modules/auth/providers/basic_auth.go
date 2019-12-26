package providers

import (
	"fmt"
	"net/http"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
)

// BasicAuthProvider ...
type BasicAuthProvider struct {
	template
	writer http.ResponseWriter
}

// NewBasicAuthProvider ...
func NewBasicAuthProvider(context context.Context, s session.Session, r *http.Request, w http.ResponseWriter) *BasicAuthProvider {
	return &BasicAuthProvider{template: *newTemplate(context, s, r), writer: w}
}

// IsAuthenticated ...
func (p *BasicAuthProvider) IsAuthenticated() bool {
	value := p.session.GetToken()
	return value == "admin"
}

// Logout ...
func (p *BasicAuthProvider) Logout() (handled bool) {
	p.request.SetBasicAuth("", "")
	return p.template.Logout()
}

// Login ...
func (p *BasicAuthProvider) Login() (isAuth bool, handled bool) {

	user, pass, ok := p.request.BasicAuth()
	if !ok {
		p.requireBasicAuth()
		return false, true
	}
	isValid := user == config.Vars.Auth.BasicAuth.Username && pass == config.Vars.Auth.BasicAuth.Password
	if !isValid {
		// retry
		p.requireBasicAuth()
		return false, true
	}

	p.session.Save("admin")
	return isValid, false
}

func (p *BasicAuthProvider) requireBasicAuth() {
	p.writer.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, config.Vars.Auth.BasicAuth.Name))
	p.writer.WriteHeader(http.StatusUnauthorized)
	p.writer.Write([]byte("401 Unauthorized\n"))
}
