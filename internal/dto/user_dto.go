package dto

import "laundry-backend/pkg/response"

// ==========================================
// 1. REQUEST DTO (Input from Client)
// ==========================================

// CreateUserRequest defines the payload for registering a new employee.
type CreateUserRequest struct {
	FullName    string `json:"full_name" binding:"required,min=3,max=150"`
	Username    string `json:"username" binding:"required,min=3,max=100,printascii,excludesall= /\\\"'<>"`
	Email       string `json:"email" binding:"required,email,max=150"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number" binding:"required,numeric,max=30"`
	Role        string `json:"role" binding:"required,oneof=owner cashier staff courier"`
}

// UpdateUserRequest defines the payload for updating an existing employee profile.
// Pointers/OmitEmpty are used to allow partial updates (PATCH-like behavior).
type UpdateUserRequest struct {
	FullName    string `json:"full_name" binding:"omitempty,min=3,max=150"`
	Username    string `json:"username" binding:"omitempty,min=3,max=100,printascii,excludesall= /\\\"'<>"`
	Email       string `json:"email" binding:"omitempty,email,max=150"`
	Password    string `json:"password" binding:"omitempty,min=8"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,numeric,max=30"`
	Role        string `json:"role" binding:"omitempty,oneof=owner cashier staff courier"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// ==========================================
// 2. RESPONSE DTO (Output to Client)
// ==========================================

// UserSummaryResponse defines the concise user data structure.
// Typically used for list views where full details are not required.
type UserSummaryResponse struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// UserDetailResponse defines the complete user profile structure.
// Used for detail views, including contact info and timestamps.
type UserDetailResponse struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	IsActive    bool   `json:"is_active"`
	LastLoginAt string `json:"last_login_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// UserListResponse acts as a container for the Service layer to return data + pagination.
// The Handler will unwrap this and put it into the standard BaseResponse.
type UserListResponse struct {
	Data []UserSummaryResponse `json:"data"`
	Meta response.MetaData     `json:"meta"`
}
