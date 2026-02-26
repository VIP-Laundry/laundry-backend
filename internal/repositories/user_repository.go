package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"laundry-backend/internal/models"
	"laundry-backend/pkg/response"
	"time"
)

// UserRepository defines the contract for user-related database operations.
type UserRepository interface {

	// Create Operations
	InsertUser(ctx context.Context, user *models.User) error

	// Read Operations
	FetchUsers(ctx context.Context, limit, offset int, search, role string, status int) ([]models.User, int64, error)
	FindByID(ctx context.Context, id int64) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)

	// Update Operations
	UpdateUser(ctx context.Context, user *models.User) error

	// Delete Operations (Soft Delete)
	DeleteUser(ctx context.Context, id int64) error

	// Validation Helpers
	IsEmailExists(ctx context.Context, email string, excludeID int64) (bool, error)
	IsPhoneExists(ctx context.Context, phone string, excludeID int64) (bool, error)
	IsUsernameExists(ctx context.Context, username string, excludeID int64) (bool, error)

	// Auth Helper
	UpdateLastLogin(ctx context.Context, id int64) error
}

// userRepository is the concrete implementation of UserRepository using sql.DB.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// --- IMPLEMENTATION ---

// InsertUser creates a new user record in the database.
func (r *userRepository) InsertUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (full_name, username, email, password_hash, role, phone_number, is_active, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query,
		user.FullName, user.Username, user.Email, user.PasswordHash,
		user.Role, user.PhoneNumber, user.IsActive, user.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.InsertUser: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("userRepo.InsertUser.LastInsertId: %w", err)
	}

	user.ID = id
	return nil
}

// FetchUsers retrieves a list of users with pagination and filtering support.
func (r *userRepository) FetchUsers(ctx context.Context, limit, offset int, search, role string, status int) ([]models.User, int64, error) {
	// [OPTIMASI] Gunakan builder sederhana
	whereClause := "WHERE 1=1"
	var args []interface{}

	if search != "" {
		whereClause += " AND (full_name LIKE ? OR username LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	if role != "" {
		whereClause += " AND role = ?"
		args = append(args, role)
	}
	if status != -1 {
		whereClause += " AND is_active = ?"
		args = append(args, status)
	}

	// Count Query
	var totalItems int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("userRepo.FetchUsers.Count: %w", err)
	}

	// Data Query
	query := fmt.Sprintf("SELECT id, full_name, username, role, is_active FROM users %s ORDER BY created_at DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("userRepo.FetchUsers.Query: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Username, &u.Role, &u.IsActive); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, totalItems, nil
}

// FindByID retrieves a single user's detailed information by ID.
func (r *userRepository) FindByID(ctx context.Context, id int64) (*models.User, error) {

	query := `SELECT id, full_name, username, email, password_hash, role, phone_number, is_active, last_login_at, created_at, updated_at FROM users WHERE id = ?`

	var u models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.FullName, &u.Username, &u.Email, &u.PasswordHash,
		&u.Role, &u.PhoneNumber, &u.IsActive, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound // [FIX] Langsung return error sentinel
		}
		return nil, fmt.Errorf("userRepo.FindByID: %w", err)
	}
	return &u, nil
}

// FindByUsername retrieves a user by their username (used for Login).
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {

	query := `SELECT id, full_name, username, password_hash, role, is_active FROM users WHERE username = ?`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.FullName, &user.Username, &user.PasswordHash, &user.Role, &user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound // [FIX] Langsung return error sentinel
		}
		return nil, fmt.Errorf("userRepo.FindByUsername: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user record.
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {

	query := `UPDATE users SET full_name=?, username=?, email=?, password_hash=?, role=?, phone_number=?, is_active=?, updated_at=? WHERE id=?`

	_, err := r.db.ExecContext(ctx, query,
		user.FullName, user.Username, user.Email, user.PasswordHash,
		user.Role, user.PhoneNumber, user.IsActive, user.UpdatedAt, user.ID,
	)

	if err != nil {
		return fmt.Errorf("userRepo.UpdateUser: %w", err)
	}
	return nil
}

// DeleteUser performs a soft delete by setting is_active to false.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {

	query := "UPDATE users SET is_active = 0 WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("userRepo.DeleteUser: %w", err)
	}
	return nil
}

// IsEmailExists checks if an email is already in use by another user.
func (r *userRepository) IsEmailExists(ctx context.Context, email string, excludeID int64) (bool, error) {

	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)"

	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userRepo.IsEmailExists: %w", err)
	}
	return exists, nil
}

// IsPhoneExists checks if a phone number is already in use by another user.
func (r *userRepository) IsPhoneExists(ctx context.Context, phone string, excludeID int64) (bool, error) {

	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = ? AND id != ?)"

	err := r.db.QueryRowContext(ctx, query, phone, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userRepo.IsPhoneExists: %w", err)
	}
	return exists, nil
}

// IsUsernameExists checks if a username is already in use by another user.
func (r *userRepository) IsUsernameExists(ctx context.Context, username string, excludeID int64) (bool, error) {

	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND id != ?)"

	err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userRepo.IsUsernameExists: %w", err)
	}
	return exists, nil
}

// UpdateLastLogin updates the last_login_at timestamp to the current time.
func (r *userRepository) UpdateLastLogin(ctx context.Context, id int64) error {

	query := "UPDATE users SET last_login_at = ? WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		// [CONSISTENCY] Bungkus error agar seragam dengan fungsi lain
		return fmt.Errorf("userRepo.UpdateLastLogin: %w", err)
	}
	return nil
}
