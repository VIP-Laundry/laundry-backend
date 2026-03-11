package dto

import "laundry-backend/pkg/response"

// ==========================================
// 1. REQUEST DTO (Input from Client)
// ==========================================

// CreateCategoryRequest untuk payload POST /categories
type CreateCategoryRequest struct {
	CategoryName string  `json:"category_name" binding:"required,min=3,max=150"`
	Description  *string `json:"description" binding:"max=255"`
}

// UpdateCategoryRequest untuk payload PUT /categories/:id
// Menggunakan pointer (*) untuk mendukung Partial Update (membedakan null vs empty string)
type UpdateCategoryRequest struct {
	CategoryName *string `json:"category_name" binding:"omitempty,min=3,max=150"`
	Description  *string `json:"description" binding:"omitempty,max=255"`
	IsActive     *bool   `json:"is_active" binding:"omitempty"`
}

// ==========================================
// 2. RESPONSE DTO (Output to Client)
// ==========================================

// CategorySummaryResponse untuk endpoint List (GET /categories)
// Sesuai API Specs: Ringkas, tanpa timestamps, tanpa description (opsional, tapi biar ringan kita hapus description)
type CategorySummaryResponse struct {
	ID           int64   `json:"id"`
	CategoryName string  `json:"category_name"`
	Description  *string `json:"description"`
	IsActive     bool    `json:"is_active"`
}

// CategoryDetailResponse untuk endpoint Detail (GET /categories/:id)
// Menampilkan data lengkap
type CategoryDetailResponse struct {
	ID           int64   `json:"id"`
	CategoryName string  `json:"category_name"`
	Description  *string `json:"description"`
	IsActive     bool    `json:"is_active"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
}

// CategoryListResponse untuk wrapper pagination
type CategoryListResponse struct {
	Data []CategorySummaryResponse `json:"data"`
	Meta response.MetaData         `json:"meta"`
}
