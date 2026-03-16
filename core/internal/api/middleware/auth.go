package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"sentinelos/core/internal/auth"
	"sentinelos/core/pkg/utils"
)

type contextKey string

const userContextKey = contextKey("user")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5001", "Missing authorization token", "")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5002", "Invalid token format", "expected Bearer token")
			return
		}

		tokenStr := parts[1]

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return auth.JwtKey(), nil
		})

		if err != nil || !token.Valid {
			utils.SendError(w, http.StatusUnauthorized, "ERR_SEC_5003", "Invalid or expired token", err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func UserFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(userContextKey).(*auth.Claims)
	return claims, ok
}
