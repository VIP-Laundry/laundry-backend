package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"laundry-backend/internal/models"
	"laundry-backend/pkg/response"
	"strings"
)

// CategoryRepository defines the contract for service category-related database operations.
type CategoryRepository interface {

	// Create Operations
	InsertCategory(ctx context.Context, category *models.ServiceCategory) error

	// Read Operations
	FindAll(ctx context.Context, limit, offset int, search, status, sortBy, sortOrder string) ([]models.ServiceCategory, int64, error)
	FindByID(ctx context.Context, id int64) (*models.ServiceCategory, error)
	FindByName(ctx context.Context, categoryName string) (*models.ServiceCategory, error)

	// Update Operations
	UpdateCategory(ctx context.Context, category *models.ServiceCategory) error

	// Delete Operations (Soft Delete)
	DeleteCategory(ctx context.Context, id int64) error
}

// categoryRepository is the concrete implementation of CategoryRepository using sql.DB.
type categoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new instance of CategoryRepository.
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// --- IMPLEMENTATION ---

// InsertCategory creates a new service category record in the database.
func (r *categoryRepository) InsertCategory(ctx context.Context, category *models.ServiceCategory) error {

	// 1. Persiapkan query SQL
	query := `
		INSERT INTO service_categories (category_name, description, is_active, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`

	// 2. Eksekusi query dengan context (Pilar O - Optimal Go)
	res, err := r.db.ExecContext(ctx, query,
		category.CategoryName,
		category.Description, // Pointer, aman jika nil
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt, // Pointer, aman jika nil
	)

	if err != nil {
		// 3. Bungkus error agar mudah dilacak (Pilar E - Error Handling)
		return fmt.Errorf("categoryRepo.InsertCategory.Exec: %w", err)
	}

	// 4. Ambil ID yang baru saja di-generate oleh MySQL (Auto Increment)
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("categoryRepo.InsertCategory.LastInsertId: %w", err)
	}

	// 5. Sematkan ID kembali ke struct pointer
	category.ID = id
	return nil
}

// FindAll retrieves a list of service categories with pagination, filtering, and sorting support.
func (r *categoryRepository) FindAll(ctx context.Context, limit, offset int, search, status, sortBy, sortOrder string) ([]models.ServiceCategory, int64, error) {

	// 1. Inisialisasi query dasar menggunakan builder sederhana
	whereClause := "WHERE 1=1"
	var args []interface{}
		
	// 2. Terapkan filter pencarian nama kategori (LOWER untuk case-insensitive)
	if search != "" {
		whereClause += " AND LOWER(category_name) LIKE ?"
		args = append(args, "%"+strings.ToLower(search)+"%")
	}

	// 3. Terapkan filter status aktif/non-aktif
	if status != "" {
		if status == "1" {
			whereClause += " AND is_active = 1"
		} else if status == "0" {
			whereClause += " AND is_active = 0"
		}
	}

	// 4. Hitung total baris keseluruhan (Count Query) untuk data Meta Pagination
	var totalItems int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM service_categories %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("categoryRepo.FindAll.Count: %w", err)
	}

	// 5. Validasi kolom sorting (Kunci keamanan mencegah SQL Injection via ORDER BY)
	validSortColumns := map[string]bool{
		"category_name": true,
		"created_at":    true,
		"id":            true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "category_name" // Fallback ke default
	}

	// 6. Validasi arah sorting (ASC/DESC)
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "ASC" // Fallback ke default
	}

	// 7. Rangkai query utama beserta LIMIT dan OFFSET
	query := fmt.Sprintf("SELECT id, category_name, description, is_active, created_at, updated_at FROM service_categories %s ORDER BY %s %s LIMIT ? OFFSET ?", whereClause, sortBy, sortOrder)
	args = append(args, limit, offset)

	// 8. Eksekusi query utama
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("categoryRepo.FindAll.Query: %w", err)
	}
	defer rows.Close()

	// 9. Parsing (Mapping) hasil query ke dalam slice struct
	var categories []models.ServiceCategory
	for rows.Next() {
		var c models.ServiceCategory

		// 10. Siapkan wadah perantara untuk menangkap NULL dari database MySQL
		var descriptionNull sql.NullString
		var updatedAtNull sql.NullTime

		if err := rows.Scan(&c.ID, &c.CategoryName, &descriptionNull, &c.IsActive, &c.CreatedAt, &updatedAtNull); err != nil {
			return nil, 0, fmt.Errorf("categoryRepo.FindAll.Scan: %w", err)
		}

		// 11. Pindahkan isi dari wadah perantara ke pointer struct Model jika datanya valid (bukan NULL)
		if descriptionNull.Valid {
			c.Description = &descriptionNull.String
		}
		if updatedAtNull.Valid {
			c.UpdatedAt = &updatedAtNull.Time
		}

		categories = append(categories, c)
	}

	return categories, totalItems, nil
}

// FindByID retrieves a single service category's detailed information by ID.
func (r *categoryRepository) FindByID(ctx context.Context, id int64) (*models.ServiceCategory, error) {

	// 1. Persiapkan query
	query := `SELECT id, category_name, description, is_active, created_at, updated_at FROM service_categories WHERE id = ?`

	var c models.ServiceCategory
	var descriptionNull sql.NullString
	var updatedAtNull sql.NullTime

	// 2. Eksekusi query dan mapping (Scan)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.CategoryName, &descriptionNull, &c.IsActive, &c.CreatedAt, &updatedAtNull,
	)

	// 3. Tangani error, kembalikan Sentinel Error standar VIP jika data kosong
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound // [FIX] Konsisten dengan Auth & User Repo
		}
		return nil, fmt.Errorf("categoryRepo.FindByID: %w", err)
	}

	// 4. Konversi nilai NULL database ke Pointer Model
	if descriptionNull.Valid {
		c.Description = &descriptionNull.String
	}
	if updatedAtNull.Valid {
		c.UpdatedAt = &updatedAtNull.Time
	}

	return &c, nil
}

// FindByName retrieves a single service category by its exact name (Useful for duplicate validation).
func (r *categoryRepository) FindByName(ctx context.Context, categoryName string) (*models.ServiceCategory, error) {

	// 1. Persiapkan query
	query := `SELECT id, category_name, description, is_active, created_at, updated_at FROM service_categories WHERE category_name = ?`

	var c models.ServiceCategory
	var descriptionNull sql.NullString
	var updatedAtNull sql.NullTime

	// 2. Eksekusi query dan mapping (Scan)
	err := r.db.QueryRowContext(ctx, query, categoryName).Scan(
		&c.ID, &c.CategoryName, &descriptionNull, &c.IsActive, &c.CreatedAt, &updatedAtNull,
	)

	// 3. Tangani error (Sentinel Error)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("categoryRepo.FindByName: %w", err)
	}

	// 4. Konversi nilai NULL database ke Pointer Model
	if descriptionNull.Valid {
		c.Description = &descriptionNull.String
	}
	if updatedAtNull.Valid {
		c.UpdatedAt = &updatedAtNull.Time
	}

	return &c, nil
}

// UpdateCategory updates an existing service category record.
func (r *categoryRepository) UpdateCategory(ctx context.Context, category *models.ServiceCategory) error {

	// 1. Persiapkan query UPDATE
	query := `UPDATE service_categories 
        SET category_name = ?, description = ?, is_active = ?, updated_at = ? 
        WHERE id = ?`

	// 2. Eksekusi query
	res, err := r.db.ExecContext(ctx, query,
		category.CategoryName,
		category.Description,
		category.IsActive,
		category.UpdatedAt,
		category.ID,
	)

	// 3. Tangani error eksekusi dan bungkus dengan rapi
	if err != nil {
		return fmt.Errorf("categoryRepo.UpdateCategory.Exec: %w", err)
	}

	// --- [CEK TAMBAHAN: SATPAM LOGIKA] ---
	// Kita cek apakah ada baris yang benar-benar berubah di database
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("categoryRepo.UpdateCategory.RowsAffected: %w", err)
	}

	if rows == 0 {
		// Jika 0, berarti ID tidak ditemukan atau data yang dikirim sama persis dengan yang lama
		return response.ErrNotFound
	}

	return nil
}

// DeleteCategory performs a soft delete by setting is_active to false (0).
func (r *categoryRepository) DeleteCategory(ctx context.Context, id int64) error {

	// 1. Persiapkan query soft-delete
	query := "UPDATE service_categories SET is_active = 0 WHERE id = ? AND is_active = 1"

	// 2. EKSEKUSI (Jalankan)
	// Kita tangkap hasilnya di variabel 'res'
	res, err := r.db.ExecContext(ctx, query, id)

	// 3. EVALUASI (Cek Error Teknis)
	if err != nil {
		return fmt.Errorf("categoryRepo.DeleteCategory.Exec: %w", err)
	}

	// --- [SATPAM TAMBAHAN: RowsAffected] ---
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("categoryRepo.DeleteCategory.RowsAffected: %w", err)
	}

	if rows == 0 {
		// Jika 0, berarti ID tidak ditemukan atau kategori sudah non-aktif sebelumnya
		return response.ErrNotFound
	}

	// 4. KONKLUSI (Kembalikan)
	return nil
}
