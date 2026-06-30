package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/hazubeep/personal-blog/internal/config"
)

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	cfg := config.Load()

	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok ||
			subtle.ConstantTimeCompare([]byte(username), []byte(cfg.AdminUser)) != 1 ||
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.AdminPass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted area"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
