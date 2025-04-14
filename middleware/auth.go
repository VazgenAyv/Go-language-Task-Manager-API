package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ht21992/go-task-manager/models"
)

var jwtKey = []byte("my_secret_key")

type contextKey string

const (
	UserContextKey = contextKey("user")
	RoleContextKey = contextKey("role")
)

// JWTMiddleware verifies the JWT token and adds user info to the request context
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims.Username)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)

		fmt.Printf("Claims: %+v\n", *claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
