package providers

import (
	"encoding/json"
	"fmt"
	"github.com/tstranex/u2f"
	"html/template"
	"net/http"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
	"secure-portal/modules/util"
)

// FIDOProvider ...
type FIDOProvider struct {
	providerTemplate
	counter       uint32
	writer        http.ResponseWriter
	challenge     *u2f.Challenge
	registrations []u2f.Registration

	templateLogin            *template.Template
	templateRegister         *template.Template
	templateLoginValidate    *template.Template
	templateRegisterValidate *template.Template
}

// NewFIDOProvider ...
func NewFIDOProvider(context context.Context, s session.Session, r *http.Request, w http.ResponseWriter) *FIDOProvider {
	provider := &FIDOProvider{providerTemplate: *newTemplate(context, s, r), writer: w}
	provider.init()
	return provider
}

// IsAuthenticated ...
func (p *FIDOProvider) init() {
	var err error
	dir := "/templates/auth/fido"

	//TODO move this load template to another package
	p.templateLogin, err = util.LoadTemplate(fmt.Sprintf("%s/%s", dir, "login.html"))
	if err != nil {
		panic(err)
	}
	p.templateLoginValidate, err = util.LoadTemplate(fmt.Sprintf("%s/%s", dir, "login_validate.html"))
	if err != nil {
		panic(err)
	}
	p.templateRegister, err = util.LoadTemplate(fmt.Sprintf("%s/%s", dir, "register.html"))
	if err != nil {
		panic(err)
	}
	p.templateRegisterValidate, err = util.LoadTemplate(fmt.Sprintf("%s/%s", dir, "register_validate.html"))
	if err != nil {
		panic(err)
	}
}

// IsAuthenticated ...
func (p *FIDOProvider) IsAuthenticated() bool {
	value := p.session.GetToken()
	return value == "admin"
}

// Logout ...
func (p *FIDOProvider) Logout() (handled bool) {
	p.request.SetBasicAuth("", "")
	return p.providerTemplate.Logout()
}

// appID ...
func (p *FIDOProvider) appID() string {
	host := config.Vars.Auth.Target.Host
	if host == "" {
		host = p.request.Host
	}
	return fmt.Sprintf("%s://%s", config.Vars.Auth.Target.Schema, host)
}

// trustedFacets ...
func (p *FIDOProvider) trustedFacets() []string {
	return []string{p.appID()}
}

// Register ...
func (p *FIDOProvider) Register() (handled bool) {

	if p.request.Method == http.MethodPost {
		return p.RegisterValidate()
	}

	log := p.Ctx.GetLogger("Register")
	c, err := u2f.NewChallenge(p.appID(), p.trustedFacets())
	if err != nil {
		log.Warnf("challenge error: %s", err.Error())
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return true
	}

	req := u2f.NewWebRegisterRequest(c, p.registrations)
	p.challenge = c

	log.Infof("registerRequest %+v", req)

	params := map[string]string{}

	request, err := json.Marshal(req)
	params["request"] = string(request)
	params["registerPath"] = config.Vars.Auth.Path.Register

	err = p.templateRegister.Execute(p.writer, params)
	if err != nil {
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return true
	}

	return true
}

// RegisterValidate ...
func (p *FIDOProvider) RegisterValidate() (handled bool) {

	log := p.Ctx.GetLogger("RegisterValidate")

	var regResp u2f.RegisterResponse
	if err := json.NewDecoder(p.request.Body).Decode(&regResp); err != nil {
		http.Error(p.writer, "invalid response: "+err.Error(), http.StatusBadRequest)
		return true
	}

	if p.challenge != nil {
		http.Error(p.writer, "challenge missing", http.StatusBadRequest)
		return true
	}

	challengeConfig := &u2f.Config{
		// Chrome 66+ doesn't return the device's attestation
		// certificate by default.
		SkipAttestationVerify: true,
	}

	reg, err := u2f.Register(regResp, *p.challenge, challengeConfig)
	if err != nil {
		log.Printf("u2f.Register error: %v", err)
		http.Error(p.writer, "error verifying response", http.StatusInternalServerError)
		return true
	}

	p.registrations = append(p.registrations, *reg)
	p.counter = 0

	log.Printf("Registration success: %+v", reg)
	p.writer.Write([]byte("success"))
	return true
}

// Login ...
func (p *FIDOProvider) Login() (isAuth bool, handled bool) {

	if p.request.Method == http.MethodPost {
		return p.LoginValidate()
	}

	log := p.Ctx.GetLogger("Login")

	if p.challenge != nil {
		http.Error(p.writer, "challenge missing", http.StatusBadRequest)
		return false, true
	}

	if p.registrations == nil {
		http.Error(p.writer, "registration missing", http.StatusBadRequest)
		return false, true
	}

	c, err := u2f.NewChallenge(p.appID(), p.trustedFacets())
	if err != nil {
		log.Warnf("challenge error: %s", err.Error())
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return false, true
	}
	p.challenge = c

	req := c.SignRequest(p.registrations)

	log.Printf("Sign request: %+v", req)

	params := map[string]string{}

	request, err := json.Marshal(req)
	params["request"] = string(request)
	params["loginPath"] = config.Vars.Auth.Path.Login

	err = p.templateLogin.Execute(p.writer, params)
	if err != nil {
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return false, true
	}

	return false, true
}

// LoginValidate ...
func (p *FIDOProvider) LoginValidate() (isAuth bool, handled bool) {

	log := p.Ctx.GetLogger("LoginValidate")

	var signResp u2f.SignResponse
	if err := json.NewDecoder(p.request.Body).Decode(&signResp); err != nil {
		http.Error(p.writer, "invalid response: "+err.Error(), http.StatusBadRequest)
		return false, true
	}

	log.Printf("signResponse: %+v", signResp)

	if p.registrations == nil {
		http.Error(p.writer, "registration missing", http.StatusBadRequest)
		return false, true
	}

	var err error
	for _, reg := range p.registrations {
		newCounter, authErr := reg.Authenticate(signResp, *p.challenge, p.counter)
		if authErr == nil {
			log.Printf("newCounter: %d", newCounter)
			p.counter = newCounter
			return true, false
		}
	}

	log.Printf("VerifySignResponse error: %v", err)
	http.Error(p.writer, "error verifying response", http.StatusInternalServerError)

	return false, true
}
