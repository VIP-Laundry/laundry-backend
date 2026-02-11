package repositories

import (
	"context"
	"database/sql"
	"errors"
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
		INSERT INTO users (full_name, username, email, password_hash, role, phone_number, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(ctx, query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PhoneNumber,
		user.IsActive,
		user.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

// FetchUsers retrieves a list of users with pagination and filtering support.
func (r *userRepository) FetchUsers(ctx context.Context, limit, offset int, search, role string, status int) ([]models.User, int64, error) {

	// 1. Base Query Construction
	query := "SELECT id, full_name, username, role, is_active FROM users WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	var args []interface{}

	// 2. Dynamic Filtering
	if search != "" {
		filter := " AND (full_name LIKE ? OR username LIKE ?)"
		query += filter
		countQuery += filter
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if role != "" {
		filter := " AND role = ?"
		query += filter
		countQuery += filter
		args = append(args, role)
	}

	if status != -1 { // -1 means ignore status filter
		filter := " AND is_active = ?"
		query += filter
		countQuery += filter
		args = append(args, status)
	}

	// 3. Count Total Items (for Pagination Meta)
	var totalItems int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// 4. Pagination & Sorting
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	// 5. Execute Data Query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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

	query := `
		SELECT id, full_name, username, email, password_hash, role, phone_number, is_active, last_login_at, created_at, updated_at 
		FROM users WHERE id = ?
	`

	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.FullName,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.PhoneNumber,
		&u.IsActive,
		&u.LastLoginAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(response.ErrNotFound)
		}
		return nil, err
	}
	return &u, nil
}

// FindByUsername retrieves a user by their username (used for Login).
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, full_name, username, password_hash, role, is_active 
		FROM users WHERE username = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.FullName,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(response.ErrNotFound)
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record.
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET full_name=?, username=?, email=?, password_hash=?, role=?, phone_number=?, is_active=?, updated_at=? 
		WHERE id=?
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FullName,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.PhoneNumber,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

// DeleteUser performs a soft delete by setting is_active to false.
func (r *userRepository) DeleteUser(ctx context.Context, id int64) error {
	query := "UPDATE users SET is_active = 0 WHERE id = ?"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IsEmailExists checks if an email is already in use by another user.
func (r *userRepository) IsEmailExists(ctx context.Context, email string, excludeID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND id != ?)"
	err := r.db.QueryRowContext(ctx, query, email, excludeID).Scan(&exists)
	return exists, err
}

// IsPhoneExists checks if a phone number is already in use by another user.
func (r *userRepository) IsPhoneExists(ctx context.Context, phone string, excludeID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE phone_number = ? AND id != ?)"
	err := r.db.QueryRowContext(ctx, query, phone, excludeID).Scan(&exists)
	return exists, err
}

// IsUsernameExists checks if a username is already in use by another user.
func (r *userRepository) IsUsernameExists(ctx context.Context, username string, excludeID int64) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND id != ?)"
	err := r.db.QueryRowContext(ctx, query, username, excludeID).Scan(&exists)
	return exists, err
}

// UpdateLastLogin updates the last_login_at timestamp to the current time.
func (r *userRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	query := "UPDATE users SET last_login_at = ? WHERE id = ?"
	// used time.Now() in here
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
