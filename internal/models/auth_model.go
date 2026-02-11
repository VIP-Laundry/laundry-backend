package models

import "time"

// RefreshToken represents the data structure for handling user session renewal.
// It maps to the 'refresh_tokens' table in the database.
type RefreshToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"` // UUID string
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenBlacklist stores JTI (JWT ID) of invalidated tokens (e.g., during logout).
// It maps to the 'token_blacklist' table.
type TokenBlacklist struct {
	ID        int64     `json:"id"`
	JTI       string    `json:"jti"` // Uses JTI (all caps) to match JWT standard
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
