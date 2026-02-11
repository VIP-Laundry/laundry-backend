package repositories

import (
	"context"
	"database/sql"
	"errors"
	"laundry-backend/internal/models"
	"laundry-backend/pkg/response"
)

// AuthRepository defines the contract for authentication-related database interactions.
type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken *models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, token string) error
	GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error)
	AddToBlacklist(ctx context.Context, blacklist *models.TokenBlacklist) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

// authRepository is a concrete implementation of the AuthRepository interface, which uses sql.DB as its database engine.
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateRefreshToken persists a new refresh token into the database.
func (r *authRepository) CreateRefreshToken(ctx context.Context, rt *models.RefreshToken) error {
	query := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, rt.UserID, rt.Token, rt.ExpiresAt)
	return err
}

// AddToBlacklist inserts a JTI (JWT ID) into the blacklist table.
func (r *authRepository) AddToBlacklist(ctx context.Context, tb *models.TokenBlacklist) error {
	query := "INSERT INTO token_blacklist (jti, expires_at) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, tb.JTI, tb.ExpiresAt)
	return err
}

// IsBlacklisted checks if a JTI is present in the blacklist.
func (r *authRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE jti = ?)"
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&exists)
	return exists, err
}

// DeleteRefreshToken removes a refresh token (used during logout).
func (r *authRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := "DELETE FROM refresh_tokens WHERE token = ?"
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// GetRefreshToken retrieves refresh token details by its token string.
func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	query := "SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = ?"
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Consistency: using the same error string as UserRepository
			return nil, errors.New(response.ErrNotFound)
		}
		return nil, err
	}
	return &rt, nil
}
