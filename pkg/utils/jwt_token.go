package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT payload data.
type Claims struct {
	UserID   int64  `json:"user_id"` // Explicit int64 prevents float64 unmarshaling panic
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new HS256 JWT. Returns the token string and JTI.
func GenerateAccessToken(userID int64, username, role string, secretKey []byte, expiry time.Duration) (string, string, error) {
	jti := uuid.New().String()
	expirationTime := time.Now().Add(expiry)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vip-laundry-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)

	return tokenString, jti, err
}

// ValidateAccessToken parses and validates the token string.
func ValidateAccessToken(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

// GenerateRefreshToken generates a secure random UUID for session renewal.
func GenerateRefreshToken() (string, error) {
	return uuid.New().String(), nil
}
