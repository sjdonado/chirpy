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
		return "", fmt.Errorf("Missing Bearer token")
	}

	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := strings.TrimSpace(headers.Get("Authorization"))

	token := strings.TrimPrefix(authHeader, "ApiKey ")
	if token == "" {
		return "", fmt.Errorf("Missing ApiKey token")
	}

	return token, nil
}
