package middleware

import "net/http"

func SecurityHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-XSS-Protection", "0") // modern browsers; rely on CSP later if needed
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		// Only enable HSTS when you are fully on HTTPS in production
		// w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

		next.ServeHTTP(w, r)
	})
}
