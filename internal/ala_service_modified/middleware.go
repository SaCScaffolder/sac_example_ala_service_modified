package ala_service_modified

import "net/http"

// CORSHandler wraps next with permissive CORS headers. Exported as a
// package-level helper so cmd/server/main.go can compose it without
// referencing the lowercase function name.
func CORSHandler(next http.Handler) http.Handler {
	return corsMiddleware(next)
}

// corsMiddleware sets the standard CORS headers on every response and
// short-circuits preflight OPTIONS requests. It's intentionally
// permissive (AllowOrigins=*) for the integration-test fixture;
// production services should restrict AllowOrigins to the real
// frontend origin.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}