package middleware

import (
	"context"
	"net/http"
	"strings"

	"folder-system/internal/utils"
)

func AuthMiddleware(accessSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Authorization header is required"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "Authorization header format must be Bearer {token}"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := utils.ParseJWT(tokenString, accessSecret)
			if err != nil {
				http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Add user ID to the request context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
