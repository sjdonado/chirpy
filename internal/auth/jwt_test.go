package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	expiresIn   = 24 * time.Hour
	tokenSecret = "secret123"
)

var userId = uuid.New()

// TestMakeJWT validates that MakeJWT returns a valid JWT token.
func TestMakeJWT(t *testing.T) {
	jwtToken, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT(%v, %s, %v) returned %v`, userId, tokenSecret, expiresIn, err)
		return
	}

	// validate jwt is valid base64 string
	if _, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}); err != nil {
		t.Errorf(`MakeJWT(%v, %s, %v) returned invalid JWT token`, userId, tokenSecret, expiresIn)
		return
	}
}

// TestValidateJWT validates that ValidateJWT returns a valid user ID.
func TestValidateJWT(t *testing.T) {
	jwtToken, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf(`MakeJWT(%v, %s, %v) returned %v`, userId, tokenSecret, expiresIn, err)
		return
	}

	userID, err := ValidateJWT(jwtToken, tokenSecret)
	if err != nil {
		t.Errorf(`ValidateJWT(%v, %s) returned %v`, jwtToken, tokenSecret, err)
		return
	}

	if userID != userId {
		t.Errorf(`ValidateJWT(%v, %s) returned %v, want %v`, jwtToken, tokenSecret, userID, userId)
		return
	}
}
