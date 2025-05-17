package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecretKey is used to sign and verify JWT tokens.
// ⚠️ In production, this should be loaded from a secure environment variable, not hardcoded.
var SecretKey = []byte("super-secret") // TODO: replace with os.Getenv("JWT_SECRET")

// Generate creates a new JWT token with an email claim and 24h expiration.
// The token can be used for confirming subscriptions or unsubscribing without login.
func Generate(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// Parse validates the JWT token and extracts the email claim.
// Returns an error if the token is invalid or the claim is missing/malformed.
func Parse(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	// Safely extract the "email" claim
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrTokenMalformed
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", jwt.ErrTokenMalformed
	}

	return email, nil
}
