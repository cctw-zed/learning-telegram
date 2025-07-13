package api

import (
	"context"
	"net/http"
	"strings"

	"learning-telegram/internal/auth"
)

// CORS中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware 验证JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "需要认证", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "无效的认证格式", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "无效的token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
