package middleware

import (
	"crypto/subtle"
	"net/http"
)

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok ||
			subtle.ConstantTimeCompare([]byte(username), []byte("admin")) != 1 ||
			subtle.ConstantTimeCompare([]byte(password), []byte("admin")) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted area"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
