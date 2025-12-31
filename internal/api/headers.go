package api

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := strings.TrimSpace(headers.Get("Authorization"))

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("Missing bearer token")
	}

	return token, nil
}
