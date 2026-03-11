package services

import (
	"context"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/response"
	"time"
)

// CategoryService defines the contract for business logic related to service categories.
type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryDetailResponse, error)
	GetCategoryList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.CategoryListResponse, error)
	GetCategoryDetail(ctx context.Context, id int64) (*dto.CategoryDetailResponse, error)

	// ModifyCategoryData updates category information with validation logic.
	ModifyCategory(ctx context.Context, targetID int64, req dto.UpdateCategoryRequest) (*dto.CategoryDetailResponse, error)

	// DeleteCategory handles soft deletion of a category.
	DeactivateCategory(ctx context.Context, targetID int64) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

// NewCategoryService creates a new instance of CategoryService.
func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

// CreateCategory handles the creation of a new service classification.
func (s *categoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryDetailResponse, error) {

	// 1. Pengecekan Duplikasi Data (Nama harus unik)
	existingCategory, _ := s.categoryRepo.FindByName(ctx, req.CategoryName)
	if existingCategory != nil {
		return nil, response.ErrDuplicate
	}

	// 2. Siapkan Model (Wadah untuk dikirim ke Database)
	categoryModel := &models.ServiceCategory{
		CategoryName: req.CategoryName,
		Description:  req.Description,
		IsActive:     true, // Default aktif saat pertama dibuat
		CreatedAt:    time.Now(),
		// UpdatedAt tidak perlu ditulis nil, karena otomatis nil bawaan Go
	}

	// 3. Insert ke Database (Pointer Magic bekerja di sini)
	err := s.categoryRepo.InsertCategory(ctx, categoryModel)
	if err != nil {
		return nil, err
	}

	// 4. Kembalikan Response menggunakan Fungsi Helper
	return s.mapToDetailResponse(categoryModel), nil

}

// GetCategoryList fetches a list of categories with pagination, filters, and sorting.
func (s *categoryService) GetCategoryList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.CategoryListResponse, error) {

	// 1. Validasi Batas Halaman (Safety check)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// --- [PERBAIKAN KUNCI 1: Konversi Page ke Offset] ---
	// Repository kita (SQL) butuhnya 'limit' dan 'offset', bukan 'page'
	limit := perPage
	offset := (page - 1) * perPage

	// 2. Panggil Repository (Gunakan nama fungsi yang benar)
	// --- [PERBAIKAN KUNCI 2: FetchCategories -> FindAll] ---
	categories, totalItems, err := s.categoryRepo.FindAll(ctx, limit, offset, search, status, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// 3. Mapping dari Model (Database) ke DTO Summary (Respon JSON)
	var categoryResponses []dto.CategorySummaryResponse
	for _, c := range categories {
		categoryResponses = append(categoryResponses, dto.CategorySummaryResponse{
			ID:           c.ID,
			CategoryName: c.CategoryName,
			Description:  c.Description,
			IsActive:     c.IsActive,
		})
	}

	// Cegah nilai "null" di JSON jika database kosong (Membuatnya jadi "[]")
	if categoryResponses == nil {
		categoryResponses = []dto.CategorySummaryResponse{}
	}

	// 4. Hitung Total Halaman (Total Pages)
	// Using integer math to avoid float precision issues
	var totalPages int
	if perPage > 0 {
		totalPages = int((totalItems + int64(perPage) - 1) / int64(perPage))
	}

	// 5. Kembalikan Response Akhir beserta Meta Data
	return &dto.CategoryListResponse{
		Data: categoryResponses,
		Meta: response.MetaData{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

// GetCategoryDetail retrieves detailed category information by ID.
func (s *categoryService) GetCategoryDetail(ctx context.Context, id int64) (*dto.CategoryDetailResponse, error) {

	// 1. Ambil Data dari Repository (Gunakan nama fungsi VIP-7: FindByID)
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		// Jika ID tidak ada, Repository kita sudah pintar dan akan otomatis
		// mengembalikan error standar: response.ErrNotFound
		return nil, err
	}

	// 2. Kembalikan Response menggunakan "Mesin Packing" (Helper)
	return s.mapToDetailResponse(category), nil
}

// ModifyCategoryData updates category profile with validation logic.
func (s *categoryService) ModifyCategory(ctx context.Context, targetID int64, req dto.UpdateCategoryRequest) (*dto.CategoryDetailResponse, error) {

	// 1. Ambil Data Kategori yang Lama (Gunakan nama fungsi VIP-7: FindByID)
	existingCategory, err := s.categoryRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// 2. Update Fields (Partial Update Logic)

	// Validasi Nama: Cek apakah user mengirim data nama (tidak nil)
	// Lalu buka isinya pakai *req.CategoryName untuk dicek dan di-update
	if req.CategoryName != nil && *req.CategoryName != existingCategory.CategoryName {
		duplicateCheck, _ := s.categoryRepo.FindByName(ctx, *req.CategoryName)

		// Jika nama sudah dipakai oleh ID Kategori lain, Tolak!
		if duplicateCheck != nil && duplicateCheck.ID != targetID {
			return nil, response.ErrDuplicate
		}
		existingCategory.CategoryName = *req.CategoryName
	}

	// Update Deskripsi jika dikirim oleh user (tidak nil)
	if req.Description != nil {
		// [FIX] Karena existingCategory.Description juga sebuah Pointer (*string),
		// kita langsung saja oper req.Description (yang juga Pointer) tanpa tanda bintang.
		existingCategory.Description = req.Description
	}

	// Boolean Check (Ini sudah benar, karena req.IsActive juga pointer)
	if req.IsActive != nil {
		existingCategory.IsActive = *req.IsActive
	}

	// 3. Update Waktu (Timestamp)
	now := time.Now()
	existingCategory.UpdatedAt = &now

	// 4. Simpan Perubahan ke Database
	// [FIX] UpdateCategory di VIP-7 hanya mengembalikan error, karena existingCategory sudah terupdate otomatis via Pointer
	err = s.categoryRepo.UpdateCategory(ctx, existingCategory)
	if err != nil {
		return nil, err
	}

	// 5. Kembalikan Response menggunakan "Mesin Packing" (Helper)
	// [FIX] Menggunakan mapToDetailResponse agar ringkas dan DRY (Don't Repeat Yourself)
	return s.mapToDetailResponse(existingCategory), nil
}

// DeleteCategory handles soft deletion of a category.
func (s *categoryService) DeactivateCategory(ctx context.Context, targetID int64) error {

	// 1. Cek apakah kategori tersebut ada di database
	// [FIX] Gunakan FindByID (sesuai nama fungsi di Repository baru)
	_, err := s.categoryRepo.FindByID(ctx, targetID)
	if err != nil {
		// Jika tidak ditemukan, Repository otomatis mengirim response.ErrNotFound
		return err
	}

	// 2. Eksekusi Soft Delete (Mengubah is_active menjadi false)
	// [FIX] Gunakan Delete (sesuai nama fungsi di Repository baru)
	err = s.categoryRepo.DeleteCategory(ctx, targetID)
	if err != nil {
		return err
	}

	return nil
}

// --- HELPER FUNCTION ---

func (s *categoryService) mapToDetailResponse(category *models.ServiceCategory) *dto.CategoryDetailResponse {

	createdAtStr := category.CreatedAt.Format("2006-01-02 15:04:05")
	var updatedAtPtr *string
	if category.UpdatedAt != nil {
		formatted := category.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtPtr = &formatted
	}

	return &dto.CategoryDetailResponse{
		ID:           category.ID,
		CategoryName: category.CategoryName,
		Description:  category.Description,
		IsActive:     category.IsActive,
		CreatedAt:    createdAtStr,
		UpdatedAt:    updatedAtPtr,
	}
}
