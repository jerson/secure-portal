package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secure-portal/modules/auth"
	"secure-portal/modules/auth/providers"
	"secure-portal/modules/auth/session"
	"secure-portal/modules/config"
	"secure-portal/modules/context"
)

func init() {
	err := config.ReadDefault()
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx := context.NewContextSingle("main")
	defer ctx.Close()

	log := ctx.GetLogger("main")

	origin, err := url.Parse(config.Vars.Auth.Source.Host)
	if err != nil {
		panic(err)
	}

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.NewContextSingle("request")
		defer ctx.Close()

		log := ctx.GetLogger("handler")
		provider := Provider(ctx, session.NewCookiesSession(ctx, r, w), r, w)

		log.Debugf("auth.Handler: %s", r.RequestURI)
		handled := auth.Handler(provider, w, r)
		if handled {
			log.Debug("handled")
			return
		}

		isFirstLoad := provider.IsFirstTime()
		isValid := provider.IsAuthenticated()

		log.Debug("isFirstLoad: ", isFirstLoad)
		log.Debug("isValid: ", isValid)

		if !isValid {
			if isFirstLoad {
				redirect(w, r)
				return
			}
			notAllowed(w)
			return
		}

		proxy.ServeHTTP(w, r)
	})

	port := fmt.Sprintf(":%d", config.Vars.Auth.Port)
	log.Infof("running: %s", port)

	log.Fatal(http.ListenAndServe(port, handler))

}

// Provider ...
func Provider(ctx context.Context, s session.Session, r *http.Request, w http.ResponseWriter) (provider providers.AuthProvider) {

	log := ctx.GetLogger("provider")
	log.Debugf("start: %s", config.Vars.Auth.Type)

	switch config.Vars.Auth.Type {
	case "basicauth":
		provider = providers.NewBasicAuthProvider(ctx, s, r, w)
	default:
		panic(fmt.Sprintf("type not handled: %s", config.Vars.Auth.Type))
	}

	return provider
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.Vars.Auth.Path.Login, http.StatusTemporaryRedirect)
}

func notAllowed(w http.ResponseWriter) {
	/**
	here cant use 401 header because is conflicting with auth basic
	*/
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not allowed"))
}
