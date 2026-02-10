package utils

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT payload data.
type Claims struct {
	UserID   int64  `json:"user_id"` // Explicit int64 prevents float64 unmarshaling panic
	Username string `json:"username"`
	Role     string `json:"role"`
	// JTI      string `json:"jti"`
	// REFACTOR: Field 'JTI' dihapus karena duplikat.
	// Kita akan menggunakan field bawaan 'ID' di dalam jwt.RegisteredClaims.
	jwt.RegisteredClaims
}

// GetSecretKey retrieves the secret for signing tokens from environment variables.
func GetSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// SECURITY: Memberi peringatan jika lupa set .env (Sangat berbahaya di Production)
		log.Println("⚠️ WARNING: JWT_SECRET not found in .env, using insecure default key!")
		return []byte("default_secure_secret_for_dev_only")
	}
	return []byte(secret)
}

// GetExpiryDuration retrieves the access token TTL from env or defaults to 15 minutes.
func GetExpiryDuration() time.Duration {
	expiryStr := os.Getenv("JWT_EXPIRY_MINUTE")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil || expiry <= 0 {
		return 15 * time.Minute
	}
	return time.Duration(expiry) * time.Minute
}

// GenerateAccessToken creates a new HS256 JWT. Returns the token string and JTI.
func GenerateAccessToken(userID int64, username, role string) (string, string, error) {
	jti := uuid.New().String()
	expirationTime := time.Now().Add(GetExpiryDuration())

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		// JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vip-laundry-backend", // OPTIONAL: Identitas pengeluar token
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetSecretKey())

	return tokenString, jti, err
}

// ValidateAccessToken parses and validates the token string.
func ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return GetSecretKey(), nil
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
