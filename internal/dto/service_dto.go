package dto

import (
	"laundry-backend/pkg/response"
)

// ==========================================
// REQUEST DTO (Data yang masuk dari Frontend)
// ==========================================

// CreateServiceRequest digunakan saat Owner menambah layanan baru (POST /services)
type CreateServiceRequest struct {
	CategoryID    int64   `json:"category_id" binding:"required"`
	Code          string  `json:"code" binding:"required"`
	ServiceName   string  `json:"service_name" binding:"required"`
	Unit          string  `json:"unit" binding:"required,oneof=kg pcs"`
	Price         float64 `json:"price" binding:"required,min=0"`
	DurationHours int     `json:"duration_hours" binding:"required,min=1"`
}

// UpdateServiceRequest menggunakan pointer (*) untuk mendukung Partial Update (PUT)
type UpdateServiceRequest struct {
	CategoryID    *int64   `json:"category_id"`
	Code          *string  `json:"code"`
	ServiceName   *string  `json:"service_name"`
	Unit          *string  `json:"unit" binding:"omitempty,oneof=kg pcs"`
	Price         *float64 `json:"price" binding:"omitempty,min=0"`
	DurationHours *int     `json:"duration_hours" binding:"omitempty,min=1"`
	IsActive      *bool    `json:"is_active"` // Menggunakan *bool agar bisa mendeteksi jika user mengirim 'false'
}

// ==========================================
// RESPONSE DTO (Data yang keluar ke Frontend)
// ==========================================

// 1. ServiceSummaryResponse untuk endpoint List (GET /services)
// Sesuai API Specs: Ringkas, tanpa created_at dan updated_at
type ServiceSummaryResponse struct {
	ID            int64                   `json:"id"`
	Code          string                  `json:"code"`
	ServiceName   string                  `json:"service_name"`
	Unit          string                  `json:"unit"`
	Price         float64                 `json:"price"`
	DurationHours int                     `json:"duration_hours"`
	IsActive      bool                    `json:"is_active"`
	Category      *NestedCategoryResponse `json:"category"`
}

// 2. ServiceDetailResponse untuk endpoint Detail (GET /services/:id)
// Menampilkan data lengkap beserta waktu (timestamps)
type ServiceDetailResponse struct {
	ID            int64                   `json:"id"`
	Code          string                  `json:"code"`
	ServiceName   string                  `json:"service_name"`
	Unit          string                  `json:"unit"`
	Price         float64                 `json:"price"`
	DurationHours int                     `json:"duration_hours"`
	IsActive      bool                    `json:"is_active"`
	CreatedAt     string                  `json:"created_at"`
	UpdatedAt     *string                 `json:"updated_at"` // TANPA omitempty, agar jika nil, JSON tetap mencetak "null"
	Category      *NestedCategoryResponse `json:"category"`
}

// NestedCategoryResponse untuk menyisipkan info kategori secara ringkas
type NestedCategoryResponse struct {
	ID                  int64   `json:"id"`
	CategoryName        string  `json:"category_name"`
	CategoryDescription *string `json:"category_description,omitempty"` // omitempty: sembunyikan jika nil
}

// ServiceListResponse untuk balasan GET List lengkap dengan Pagination
type ServiceListResponse struct {
	Data []ServiceSummaryResponse `json:"data"`
	Meta response.MetaData        `json:"meta"`
}
