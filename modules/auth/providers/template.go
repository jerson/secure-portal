package providers

import (
	"net/http"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/context"
)

// template ...
type template struct {
	Ctx  context.Context
	data map[string]interface{}

	session session.Session
	request *http.Request
}

// newTemplate ...
func newTemplate(ctx context.Context, s session.Session, r *http.Request) *template {
	template := &template{Ctx: ctx, session: s, request: r}
	template.init()
	return template
}

// init ...
func (p *template) init() {
	p.data = map[string]interface{}{}
}

// Session ...
func (p *template) Session() session.Session {
	return p.session
}

// IsFirstTime ...
func (p *template) IsFirstTime() bool {
	path := p.session.RedirectPath()
	if path == "" {
		p.session.SetRedirectPath(p.request.RequestURI)
		return true
	}
	return false
}

// Logout ...
func (p *template) Logout() (handled bool) {
	p.session.Reset()
	return false
}

// Register ...
func (p *template) Register() (handled bool) {
	return false
}
