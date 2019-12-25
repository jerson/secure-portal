package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
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

		if r.RequestURI == "/auth" {
			user, pass, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 Unauthorized\n"))
				return
			}
			isValidCredentials := user == "admin" && pass == "admin"
			if !isValidCredentials {
				// retry
				w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 Unauthorized\n"))
				return
			}
			if isValidCredentials {

				cookie := &http.Cookie{
					Name:  "Auth-Portal",
					Value: "admin",
					Path:  "/",
				}
				http.SetCookie(w, cookie)

				defaultPath := "/"
				defaultCookie, _ := r.Cookie("Default")
				if defaultCookie != nil {
					defaultPath = defaultCookie.Value
				}

				resetDefaultCookie := &http.Cookie{
					Name:    "Default",
					Value:   "",
					Path:    "/",
					Expires: time.Unix(0, 0),
				}
				http.SetCookie(w, resetDefaultCookie)

				http.Redirect(w, r, defaultPath, http.StatusTemporaryRedirect)
				return
			}
		}

		isFirstLoad := isFirstLoad(w, r)
		isValid := isValid(r)

		if !isValid {
			if isFirstLoad {
				http.Redirect(w, r, "/auth", http.StatusTemporaryRedirect)
				return
			}
			/**
			here cant use 401 header because is conflicting with auth basic
			*/
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not allowed"))
			return
		}

		proxy.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":9001", handler))

}

func isFirstLoad(w http.ResponseWriter, r *http.Request) bool {
	_, err := r.Cookie("Default")
	if err != nil {
		cookie := &http.Cookie{
			Name:  "Default",
			Value: r.RequestURI,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		return true
	}

	return false
}

func isValid(r *http.Request) bool {
	auth, err := r.Cookie("Auth-Portal")
	if err != nil {
		return false
	}
	return auth.Value == "admin"
}
