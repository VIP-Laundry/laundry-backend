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

// ServiceRepository adalah kontrak yang mendefinisikan semua operasi database untuk layanan.
// Layer Service HANYA akan berinteraksi dengan interface ini, bukan langsung ke struct.
type ServiceRepository interface {

	// Create Operations
	InsertService(ctx context.Context, service *models.Service) error

	// Read Operations
	FindAll(ctx context.Context, limit, offset int, search, status, sortBy, sortOrder string) ([]models.ServiceWithCategory, int64, error)
	FindByID(ctx context.Context, id int64) (*models.ServiceWithCategory, error)
	FindByCode(ctx context.Context, code string) (*models.ServiceWithCategory, error)
	FindByName(ctx context.Context, serviceName string) (*models.ServiceWithCategory, error)

	// Update Operations
	UpdateService(ctx context.Context, service *models.Service) error

	// Delete Operations (Soft Delete)
	DeleteService(ctx context.Context, id int64) error
}

// serviceRepository is the concrete implementation using sql.DB.
type serviceRepository struct {
	db *sql.DB
}

// NewServiceRepository creates a new instance of ServiceRepository.
func NewServiceRepository(db *sql.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

// --- IMPLEMENTATION ---

// InsertService creates a new service record in the database.
func (r *serviceRepository) InsertService(ctx context.Context, service *models.Service) error {

	// 1. Persiapkan query SQL
	query := `
		INSERT INTO services (code, service_name, unit, price, is_active, category_id, duration_hours, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// 2. Eksekusi query dengan context
	res, err := r.db.ExecContext(ctx, query,
		service.Code,
		service.ServiceName,
		service.Unit,
		service.Price,
		service.IsActive,
		service.CategoryID,
		service.DurationHours,
		service.CreatedAt,
		service.UpdatedAt, // Pointer, aman jika nil
	)

	// 3. Bungkus error agar mudah dilacak
	if err != nil {
		return fmt.Errorf("serviceRepo.InsertService.Exec: %w", err)
	}

	// 4. Ambil ID yang baru saja di-generate oleh MySQL (Auto Increment)
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("serviceRepo.InsertService.LastInsertId: %w", err)
	}

	// 5. Sematkan ID kembali ke struct Model
	service.ID = id
	return nil
}

// FindAll retrieves a list of services with pagination, filtering, and sorting support.
func (r *serviceRepository) FindAll(ctx context.Context, limit, offset int, search, status, sortBy, sortOrder string) ([]models.ServiceWithCategory, int64, error) {

	// 1. Inisialisasi query dasar
	whereClause := "WHERE 1=1"
	var args []interface{}

	// 2. Terapkan filter pencarian nama atau kode layanan (LOWER untuk case-insensitive)
	if search != "" {
		whereClause += " AND (LOWER(s.service_name) LIKE ? OR LOWER(s.code) LIKE ?)"
		searchParam := "%" + strings.ToLower(search) + "%"
		args = append(args, searchParam, searchParam)
	}

	// 3. Terapkan filter status aktif/non-aktif
	if status != "" {
		if status == "1" {
			whereClause += " AND s.is_active = 1"
		} else if status == "0" {
			whereClause += " AND s.is_active = 0"
		}
	}

	// 4. Hitung total baris keseluruhan (Count Query) untuk data Meta Pagination
	var totalItems int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM services s %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, fmt.Errorf("serviceRepo.FindAll.Count: %w", err)
	}

	// 5. Validasi kolom sorting (Kunci keamanan mencegah SQL Injection)
	validSortColumns := map[string]bool{
		"service_name": true,
		"code":         true,
		"price":        true,
		"created_at":   true,
		"id":           true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "s.created_at" // Fallback ke default
	} else {
		sortBy = "s." + sortBy // Tambahkan alias tabel agar tidak ambigu dengan join
	}

	// 6. Validasi arah sorting (ASC/DESC)
	sortOrder = strings.ToUpper(sortOrder)
	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC" // Fallback ke default DESC (terbaru)
	}

	// 7. Rangkai query utama dengan JOIN ke service_categories
	query := fmt.Sprintf(`
		SELECT 
			s.id, s.code, s.service_name, s.unit, s.price, s.is_active, s.created_at, s.updated_at, s.category_id, s.duration_hours,
			c.category_name, c.description AS category_description
		FROM services s
		LEFT JOIN service_categories c ON s.category_id = c.id
		%s 
		ORDER BY %s %s 
		LIMIT ? OFFSET ?`, whereClause, sortBy, sortOrder)

	args = append(args, limit, offset)

	// 8. Eksekusi query utama
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("serviceRepo.FindAll.Query: %w", err)
	}
	defer rows.Close()

	// 9. Parsing (Mapping) hasil query ke dalam slice struct
	var services []models.ServiceWithCategory
	for rows.Next() {
		var s models.ServiceWithCategory

		// 10. Siapkan wadah perantara untuk menangkap NULL dari database
		var categoryDescNull sql.NullString
		var updatedAtNull sql.NullTime

		err := rows.Scan(
			&s.ID, &s.Code, &s.ServiceName, &s.Unit, &s.Price, &s.IsActive, &s.CreatedAt, &updatedAtNull, &s.CategoryID, &s.DurationHours,
			&s.CategoryName, &categoryDescNull,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("serviceRepo.FindAll.Scan %w", err)
		}

		// 11. Pindahkan isi dari wadah perantara ke pointer struct Model
		if categoryDescNull.Valid {
			s.CategoryDescription = &categoryDescNull.String
		}
		if updatedAtNull.Valid {
			s.UpdatedAt = &updatedAtNull.Time
		}

		services = append(services, s)
	}

	return services, totalItems, nil

}

// FindByID retrieves a single service's detailed information by ID.
func (r *serviceRepository) FindByID(ctx context.Context, id int64) (*models.ServiceWithCategory, error) {

	// 1. Persiapkan query JOIN
	query := `
		SELECT 
			s.id, s.code, s.service_name, s.unit, s.price, s.is_active, s.created_at, s.updated_at, s.category_id, s.duration_hours,
			c.category_name, c.description AS category_description
		FROM services s
		LEFT JOIN service_categories c ON s.category_id = c.id
		WHERE s.id = ?
	`

	var s models.ServiceWithCategory
	var categoryDescNull sql.NullString
	var updatedAtNull sql.NullTime

	// 2. Eksekusi query dan mapping (Scan)
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID, &s.Code, &s.ServiceName, &s.Unit, &s.Price, &s.IsActive, &s.CreatedAt, &updatedAtNull, &s.CategoryID, &s.DurationHours,
		&s.CategoryName, &categoryDescNull,
	)

	// 3. Tangani error, kembalikan Sentinel Error jika data kosong
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("serviceRepo.FindByID: %w", err)
	}

	// 4. Konversi nilai NULL database ke Pointer Model
	if categoryDescNull.Valid {
		s.CategoryDescription = &categoryDescNull.String
	}
	if updatedAtNull.Valid {
		s.UpdatedAt = &updatedAtNull.Time
	}

	return &s, nil
}

// FindByCode retrieves a single service by its exact code (Useful for duplicate validation).
func (r *serviceRepository) FindByCode(ctx context.Context, code string) (*models.ServiceWithCategory, error) {

	// 1. Persiapkan query JOIN
	query := `
		SELECT 
			s.id, s.code, s.service_name, s.unit, s.price, s.is_active, s.created_at, s.updated_at, s.category_id, s.duration_hours,
			c.category_name, c.description AS category_description
		FROM services s
		LEFT JOIN service_categories c ON s.category_id = c.id
		WHERE s.code = ?
	`

	var s models.ServiceWithCategory
	var categoryDescNull sql.NullString
	var updatedAtNull sql.NullTime

	// 2. Eksekusi query dan mapping (Scan)
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&s.ID, &s.Code, &s.ServiceName, &s.Unit, &s.Price, &s.IsActive, &s.CreatedAt, &updatedAtNull, &s.CategoryID, &s.DurationHours,
		&s.CategoryName, &categoryDescNull,
	)

	// 3. Tangani error, kembalikan Sentinel Error jika data kosong
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("serviceRepo.FindByCode: %w", err)
	}

	// 4. Konversi nilai NULL database ke Pointer Model
	if categoryDescNull.Valid {
		s.CategoryDescription = &categoryDescNull.String
	}
	if updatedAtNull.Valid {
		s.UpdatedAt = &updatedAtNull.Time
	}

	return &s, nil
}

// FindByName retrieves a single service by its exact name (Useful for duplicate validation).
func (r *serviceRepository) FindByName(ctx context.Context, serviceName string) (*models.ServiceWithCategory, error) {

	// 1. Persiapkan query JOIN
	query := `
		SELECT 
			s.id, s.code, s.service_name, s.unit, s.price, s.is_active, s.created_at, s.updated_at, s.category_id, s.duration_hours,
			c.category_name, c.description AS category_description
		FROM services s
		LEFT JOIN service_categories c ON s.category_id = c.id
		WHERE s.service_name = ?
	`

	var s models.ServiceWithCategory
	var categoryDescNull sql.NullString
	var updatedAtNull sql.NullTime

	// 2. Eksekusi query dan mapping (Scan)
	err := r.db.QueryRowContext(ctx, query, serviceName).Scan(
		&s.ID, &s.Code, &s.ServiceName, &s.Unit, &s.Price, &s.IsActive, &s.CreatedAt, &updatedAtNull, &s.CategoryID, &s.DurationHours,
		&s.CategoryName, &categoryDescNull,
	)

	// 3. Tangani error, kembalikan Sentinel Error jika data kosong
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}
		return nil, fmt.Errorf("serviceRepo.FindByName: %w", err)
	}

	// 4. Konversi nilai NULL database ke Pointer Model
	if categoryDescNull.Valid {
		s.CategoryDescription = &categoryDescNull.String
	}
	if updatedAtNull.Valid {
		s.UpdatedAt = &updatedAtNull.Time
	}

	return &s, nil
}

// UpdateService updates an existing service record.
func (r *serviceRepository) UpdateService(ctx context.Context, service *models.Service) error {

	// 1. Persiapkan query UPDATE
	query := `
		UPDATE services 
		SET code = ?, service_name = ?, unit = ?, price = ?, is_active = ?, category_id = ?, duration_hours = ?, updated_at = ? 
		WHERE id = ?
	`

	// 2. Eksekusi query
	res, err := r.db.ExecContext(ctx, query,
		service.Code,
		service.ServiceName,
		service.Unit,
		service.Price,
		service.IsActive,
		service.CategoryID,
		service.DurationHours,
		service.UpdatedAt,
		service.ID,
	)

	// 3. Tangani error eksekusi
	if err != nil {
		return fmt.Errorf("serviceRepo.UpdateService.Exec: %w", err)
	}

	// --- [CEK TAMBAHAN: SATPAM LOGIKA] ---
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("serviceRepo.UpdateService.RowsAffected: %w", err)
	}

	if rows == 0 {
		return response.ErrNotFound
	}

	return nil
}

// DeleteService performs a soft delete by setting is_active to false (0).
func (r *serviceRepository) DeleteService(ctx context.Context, id int64) error {

	// 1. Persiapkan query soft-delete
	query := `
		UPDATE services SET is_active = 0 WHERE id = ? AND is_active = 1
	`

	// 2. EKSEKUSI (Jalankan)
	res, err := r.db.ExecContext(ctx, query, id)

	// 3. EVALUASI (Cek Error Teknis)
	if err != nil {
		return fmt.Errorf("serviceRepo.DeleteService.Exec: %w", err)
	}

	// --- [SATPAM TAMBAHAN: RowsAffected] ---
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("serviceRepo.DeleteService.RowsAffected: %w", err)
	}

	if rows == 0 {
		return response.ErrNotFound
	}

	// 4. KONKLUSI
	return nil
}
