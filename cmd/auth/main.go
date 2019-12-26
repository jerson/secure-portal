package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secure-portal/modules/auth"
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

	origin, _ := url.Parse(config.Vars.Auth.Source.Host)

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

		log.Debugf("auth.Handler: %s", r.RequestURI)
		handled := auth.Handler(ctx, w, r)
		if handled {
			log.Debug("handled")
			return
		}

		isFirstLoad := auth.IsFirstLoad(w, r)
		isValid := auth.IsValid(r)
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

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
}

func notAllowed(w http.ResponseWriter) {
	/**
	here cant use 401 header because is conflicting with auth basic
	*/
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not allowed"))
}
