package middleware

import (
	"context"
	"net/http"
	"strings"

	"folder-system/internal/handler"
	"folder-system/internal/utils"
)

func AuthMiddleware(accessSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handler.WriteJSONError(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				handler.WriteJSONError(w, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
				return
			}

			tokenString := parts[1]
			claims, err := utils.ParseJWT(tokenString, accessSecret)
			if err != nil {
				handler.WriteJSONError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
