package models

import "time"

// User represents the schema of the 'users' table in the database.
// It maps the columns to Go struct fields for application logic.
type User struct {
	ID           int64      `json:"id"`            // ID is the unique primary key of the user (Auto Increment).
	FullName     string     `json:"full_name"`     // FullName stores the complete name of the user.
	Username     string     `json:"username"`      // Username is a unique identifier used for login authentication.
	Email        string     `json:"email"`         // Email is a unique contact address, also used for login or notifications.
	PasswordHash string     `json:"-"`             // PasswordHash stores the encrypted password (bcrypt). It uses json:"-" tag to ensure it is never exposed in API responses.
	Role         string     `json:"role"`          // Role defines the authorization level (e.g., owner, cashier, staff, courier).
	PhoneNumber  string     `json:"phone_number"`  // PhoneNumber is the contact number of the user.
	IsActive     bool       `json:"is_active"`     // IsActive indicates the account status. true = Active, false = Inactive/Soft Deleted.
	LastLoginAt  *time.Time `json:"last_login_at"` // LastLoginAt records the timestamp of the last successful login. Pointer type (*time.Time) is used to handle NULL values from the database.
	CreatedAt    time.Time  `json:"created_at"`    // CreatedAt records the timestamp when the user account was created.
	UpdatedAt    *time.Time `json:"updated_at"`    // UpdatedAt records the timestamp of the last profile update. Pointer type (*time.Time) is used to handle NULL values.
}
