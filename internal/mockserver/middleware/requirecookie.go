package middleware

import (
	"net/http"
)

// Middleware version of RequireCookie
func RequireCookieMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(r.Cookies()) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"missing cookie"}`))
			return
		}
		next(w, r)
	}
}
