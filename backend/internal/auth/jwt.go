package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// In a real application, this should be loaded from a secure configuration.
var jwtSecret = []byte("a_very_secret_key_for_telegram_clone")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT for a given username.
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken validates a JWT string and returns the claims if valid.
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
