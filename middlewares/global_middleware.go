package middlewares

import "net/http"

// Sets the basic security headers
func SetSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Cache-Control", "no-store")
		rw.Header().Add("Content-Security-Policy", "frame-ancestors 'none'")
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("X-Content-Type-Options", "nosniff")
		rw.Header().Add("X-Frame-Options", "DENY")

		next.ServeHTTP(rw, r)
	})
}
