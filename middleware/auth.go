package middleware

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type authCtxKey string

const userIDKey authCtxKey = "userID"

type AuthMiddleware struct {
	queries *database.Queries
}

func NewAuthMiddleware(queries *database.Queries) *AuthMiddleware {
	return &AuthMiddleware{queries: queries}
}

func (m *AuthMiddleware) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := getBearerToken(r.Header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(bearerToken, os.Getenv("JWT_SECRET"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("User ID not found in context")
	}
	return userID, nil
}

func getBearerToken(headers http.Header) (string, error) {
	authHeader := strings.TrimSpace(headers.Get("Authorization"))

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("Missing bearer token")
	}

	return token, nil
}
