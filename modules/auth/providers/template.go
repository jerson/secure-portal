package providers

import (
	"net/http"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/context"
)

// providerTemplate ...
type providerTemplate struct {
	Ctx  context.Context
	data map[string]interface{}

	session session.Session
	request *http.Request
}

// newTemplate ...
func newTemplate(ctx context.Context, s session.Session, r *http.Request) *providerTemplate {
	template := &providerTemplate{Ctx: ctx, session: s, request: r}
	template.init()
	return template
}

// init ...
func (p *providerTemplate) init() {
	p.data = map[string]interface{}{}
}

// Session ...
func (p *providerTemplate) Session() session.Session {
	return p.session
}

// IsFirstTime ...
func (p *providerTemplate) IsFirstTime() bool {
	path := p.session.RedirectPath()
	if path == "" {
		p.session.SetRedirectPath(p.request.RequestURI)
		return true
	}
	return false
}

// Logout ...
func (p *providerTemplate) Logout() (handled bool) {
	p.session.Reset()
	return false
}

// Register ...
func (p *providerTemplate) Register() (handled bool) {
	return false
}
