package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("super-secret") // краще з .env

func Generate(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func Parse(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	email, ok := claims["email"].(string)
	if !ok {
		return "", jwt.ErrTokenMalformed
	}
	return email, nil
}
