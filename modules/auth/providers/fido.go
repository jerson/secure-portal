package providers

import (
	"encoding/json"
	"fmt"
	"github.com/tstranex/u2f"
	"net/http"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/context"
)

// FIDOProvider ...
type FIDOProvider struct {
	template
	counter       uint32
	writer        http.ResponseWriter
	challenge     *u2f.Challenge
	registrations []u2f.Registration
}

// NewFIDOProvider ...
func NewFIDOProvider(context context.Context, s session.Session, r *http.Request, w http.ResponseWriter) *FIDOProvider {
	return &FIDOProvider{template: *newTemplate(context, s, r), writer: w}
}

// IsAuthenticated ...
func (p *FIDOProvider) IsAuthenticated() bool {
	value := p.session.GetToken()
	return value == "admin"
}

// Logout ...
func (p *FIDOProvider) Logout() (handled bool) {
	p.request.SetBasicAuth("", "")
	return p.template.Logout()
}

// appID ...
func (p *FIDOProvider) appID() string {
	return fmt.Sprintf("%s://%s", p.request.URL.Scheme, p.request.URL.Host)
}

// trustedFacets ...
func (p *FIDOProvider) trustedFacets() []string {
	return []string{p.appID()}
}

// Register ...
func (p *FIDOProvider) Register() (handled bool) {

	log := p.Ctx.GetLogger("Register")
	if p.challenge != nil {
		return p.RegisterValidate()
	}

	c, err := u2f.NewChallenge(p.appID(), p.trustedFacets())
	if err != nil {
		log.Warnf("challenge error: %s", err.Error())
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return true
	}

	req := u2f.NewWebRegisterRequest(c, p.registrations)
	p.challenge = c

	log.Infof("registerRequest %+v", req)

	reqJSON, err := json.Marshal(req)
	if err != nil {
		log.Warnf("marshal error: %s", err.Error())
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return true
	}

	p.writer.WriteHeader(http.StatusOK)
	p.writer.Write(reqJSON)

	return true
}

// RegisterValidate ...
func (p *FIDOProvider) RegisterValidate() (handled bool) {

	log := p.Ctx.GetLogger("RegisterValidate")
	var regResp u2f.RegisterResponse
	if err := json.NewDecoder(p.request.Body).Decode(&regResp); err != nil {
		http.Error(p.writer, "invalid response: "+err.Error(), http.StatusBadRequest)
		return
	}

	config := &u2f.Config{
		// Chrome 66+ doesn't return the device's attestation
		// certificate by default.
		SkipAttestationVerify: true,
	}

	reg, err := u2f.Register(regResp, *p.challenge, config)
	if err != nil {
		log.Printf("u2f.Register error: %v", err)
		http.Error(p.writer, "error verifying response", http.StatusInternalServerError)
		return
	}

	p.registrations = append(p.registrations, *reg)
	p.counter = 0

	log.Printf("Registration success: %+v", reg)
	p.writer.Write([]byte("success"))
	return true
}

// Login ...
func (p *FIDOProvider) Login() (isAuth bool, handled bool) {

	log := p.Ctx.GetLogger("Login")

	if p.challenge != nil {
		return p.LoginValidate()
	}

	if p.registrations == nil {
		http.Error(p.writer, "registration missing", http.StatusBadRequest)
		return
	}

	c, err := u2f.NewChallenge(p.appID(), p.trustedFacets())
	if err != nil {
		log.Warnf("challenge error: %s", err.Error())
		http.Error(p.writer, "internal error", http.StatusInternalServerError)
		return
	}
	p.challenge = c

	req := c.SignRequest(p.registrations)

	log.Printf("Sign request: %+v", req)
	json.NewEncoder(p.writer).Encode(req)

	return false, false
}

// LoginValidate ...
func (p *FIDOProvider) LoginValidate() (isAuth bool, handled bool) {

	log := p.Ctx.GetLogger("Login")

	var signResp u2f.SignResponse
	if err := json.NewDecoder(p.request.Body).Decode(&signResp); err != nil {
		http.Error(p.writer, "invalid response: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("signResponse: %+v", signResp)

	if p.registrations == nil {
		http.Error(p.writer, "registration missing", http.StatusBadRequest)
		return
	}

	var err error
	for _, reg := range p.registrations {
		newCounter, authErr := reg.Authenticate(signResp, *p.challenge, p.counter)
		if authErr == nil {
			log.Printf("newCounter: %d", newCounter)
			p.counter = newCounter
			p.writer.Write([]byte("success"))
			return
		}
	}

	log.Printf("VerifySignResponse error: %v", err)
	http.Error(p.writer, "error verifying response", http.StatusInternalServerError)

	return false, false
}

const indexHTML = `
<!DOCTYPE html>
<html>
  <head>
    <script src="//code.jquery.com/jquery-1.11.2.min.js"></script>
    <!-- The original u2f-api.js code can be found here:
    https://github.com/google/u2f-ref-code/blob/master/u2f-gae-demo/war/js/u2f-api.js -->
    <script type="text/javascript" src="https://demo.yubico.com/js/u2f-api.js"></script>
  </head>
  <body>
    <h1>FIDO U2F Go Library Demo</h1>
    <ul>
      <li><a href="javascript:register();">Register token</a></li>
      <li><a href="javascript:sign();">Authenticate</a></li>
    </ul>
    <p>Open Chrome Developer Tools to see debug console logs.</p>
    <script>
  function serverError(data) {
    console.log(data);
    alert('Server error code ' + data.status + ': ' + data.responseText);
  }
  function checkError(resp) {
    if (!('errorCode' in resp)) {
      return false;
    }
    if (resp.errorCode === u2f.ErrorCodes['OK']) {
      return false;
    }
    var msg = 'U2F error code ' + resp.errorCode;
    for (name in u2f.ErrorCodes) {
      if (u2f.ErrorCodes[name] === resp.errorCode) {
        msg += ' (' + name + ')';
      }
    }
    if (resp.errorMessage) {
      msg += ': ' + resp.errorMessage;
    }
    console.log(msg);
    alert(msg);
    return true;
  }
  function u2fRegistered(resp) {
    console.log(resp);
    if (checkError(resp)) {
      return;
    }
    $.post('/registerResponse', JSON.stringify(resp)).success(function() {
      alert('Success');
    }).fail(serverError);
  }
  function register() {
    $.getJSON('/registerRequest').success(function(req) {
      console.log(req);
      u2f.register(req.appId, req.registerRequests, req.registeredKeys, u2fRegistered, 30);
    }).fail(serverError);
  }
  function u2fSigned(resp) {
    console.log(resp);
    if (checkError(resp)) {
      return;
    }
    $.post('/signResponse', JSON.stringify(resp)).success(function() {
      alert('Success');
    }).fail(serverError);
  }
  function sign() {
    $.getJSON('/signRequest').success(function(req) {
      console.log(req);
      u2f.sign(req.appId, req.challenge, req.registeredKeys, u2fSigned, 30);
    }).fail(serverError);
  }
    </script>
  </body>
</html>
`
