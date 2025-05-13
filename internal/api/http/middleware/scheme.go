package middleware

import (
	"fmt"
	"net/http"
)

// SchemeMiddleware creates a middleware that ensures a minimum scheme requirement
func SchemeMiddleware(currentScheme, requiredScheme string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isSchemeSupported(currentScheme, requiredScheme) {
				respondUnsupportedScheme(w, currentScheme, requiredScheme)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// isSchemeSupported checks if a feature is supported in the current scheme
func isSchemeSupported(currentScheme, requiredScheme string) bool {
	switch requiredScheme {
	case "basic":
		return true
	case "standard":
		return currentScheme == "standard" || currentScheme == "enhanced"
	case "enhanced":
		return currentScheme == "enhanced"
	default:
		return false
	}
}

// respondUnsupportedScheme sends an error response for unsupported scheme features
func respondUnsupportedScheme(w http.ResponseWriter, currentScheme, requiredScheme string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	message := fmt.Sprintf("This feature requires %s scheme, but current scheme is %s",
		requiredScheme, currentScheme)

	// Simple error response
	response := fmt.Sprintf(`{"success":false,"error":"%s"}`, message)
	w.Write([]byte(response))
}
