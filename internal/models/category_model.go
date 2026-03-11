package models

import "time"

// Service Category represents the `service categories` table in the database.
type ServiceCategory struct {
	ID           int64      `json:"id"`
	CategoryName string     `json:"category_name"`
	Description  *string    `json:"description"`
	IsActive     bool       `json:"is_active"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
