package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"secure-portal/modules/auth"
)

func main() {
	origin, _ := url.Parse("http://127.0.0.1:5800/")

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}

	proxy := &httputil.ReverseProxy{Director: director}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth.Handler(w, r)

		isFirstLoad := auth.IsFirstLoad(w, r)
		isValid := auth.IsValid(r)

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

	log.Fatal(http.ListenAndServe(":9001", handler))

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
