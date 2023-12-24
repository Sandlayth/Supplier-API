package helper

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func AdminAuthorizationMiddleware(next http.Handler) http.Handler {
	return authorizationMiddleware(next, "admin")
}

func ManagerAuthorizationMiddleware(next http.Handler) http.Handler {
	return authorizationMiddleware(next, "manager")
}

func authorizationMiddleware(next http.Handler, role string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			// Token is missing, respond with an authentication error
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify the token
		claims, _, err := VerifyToken(tokenString)
		if err != nil {
			// Token is invalid, respond with an authentication error
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !isRoleAllowedToAccess(claims.Role, role) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		r = r.WithContext(ctx)
		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// respondJSON is a helper function to respond with JSON data.
func RespondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func isRoleAllowedToAccess(role string, targetRole string) bool {
	if role == "admin" {
		return true
	}
	if role == "manager" && targetRole == "manager" {
		return true
	}
	return false
}

func extractTokenFromHeader(r *http.Request) string {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return ""
	}

	// Token is typically included in the format "Bearer <token>"
	splitHeader := strings.Split(authorizationHeader, " ")
	if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "bearer" {
		return ""
	}

	return splitHeader[1]
}
