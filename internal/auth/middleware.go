package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id" // Уникальный ключ
)

// JWTAuth — middleware для проверки JWT токена
func JWTAuth(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Authorization header required"})
				return
			}

			// 2. Проверяем формат: "Bearer <токен>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Invalid authorization format"})
				return
			}

			tokenString := parts[1]

			// 3. Проверяем токен
			_, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Invalid or expired token"})
				return
			}

			user, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "Invalid token"})
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

// GetUserID извлекает user_id из контекста
func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int64)
	return userID, ok
}
