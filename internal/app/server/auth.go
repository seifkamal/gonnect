package server

import "net/http"

type BasicAuthenticator interface {
	Authenticate(http.HandlerFunc) http.HandlerFunc
}

type basicAuthenticator struct {
	username string
	password string
}

func (a *basicAuthenticator) Authenticate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != a.username || pass != a.password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		f(w, r)
	}
}

func NewBasicAuthenticator(username, password string) *basicAuthenticator {
	return &basicAuthenticator{
		username: username,
		password: password,
	}
}
