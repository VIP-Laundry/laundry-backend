package services

import (
	"context"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/models"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/response"
	"time"
)

// ServiceService defines the contract for business logic related to services.
type ServiceService interface {
	CreateService(ctx context.Context, req dto.CreateServiceRequest) (*dto.ServiceDetailResponse, error)
	GetServiceList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.ServiceListResponse, error)
	GetServiceDetail(ctx context.Context, id int64) (*dto.ServiceDetailResponse, error)

	// ModifyService updates service information with validation logic.
	ModifyService(ctx context.Context, targetID int64, req dto.UpdateServiceRequest) (*dto.ServiceDetailResponse, error)

	// DeactivateService handles soft deletion of a service.
	DeactivateService(ctx context.Context, targetID int64) error
}

type serviceService struct {
	serviceRepo repositories.ServiceRepository
}

// NewServiceService creates a new instance of ServiceService.
func NewServiceService(serviceRepo repositories.ServiceRepository) ServiceService {
	return &serviceService{serviceRepo: serviceRepo}
}

// CreateService handles the creation of a new service.
func (s *serviceService) CreateService(ctx context.Context, req dto.CreateServiceRequest) (*dto.ServiceDetailResponse, error) {

	// 1. Pengecekan Duplikasi Kode (Harus unik)
	existingCode, _ := s.serviceRepo.FindByCode(ctx, req.Code)
	if existingCode != nil {
		return nil, response.ErrDuplicate
	}

	// 2. Pengecekan Duplikasi Nama (Harus unik)
	existingName, _ := s.serviceRepo.FindByName(ctx, req.ServiceName)
	if existingName != nil {
		return nil, response.ErrDuplicate
	}

	// 3. Siapkan Model (Wadah untuk dikirim ke Database)
	serviceModel := &models.Service{
		CategoryID:    req.CategoryID,
		Code:          req.Code,
		ServiceName:   req.ServiceName,
		Unit:          req.Unit,
		Price:         req.Price,
		DurationHours: req.DurationHours,
		IsActive:      true,
		CreatedAt:     time.Now(),
	}

	// 4. Insert ke Database
	err := s.serviceRepo.InsertService(ctx, serviceModel)
	if err != nil {
		return nil, err
	}

	// 5. Kembalikan Response Akhir
	// [NOTE] Berbeda dengan Category, di sini kita panggil ulang GetServiceDetail
	// agar kita mendapatkan data 'CategoryName' dan 'CategoryDescription' dari proses JOIN SQL.
	return s.GetServiceDetail(ctx, serviceModel.ID)
}

// GetServiceList fetches a list of services with pagination, filters, and sorting.
func (s *serviceService) GetServiceList(ctx context.Context, page, perPage int, search, status, sortBy, sortOrder string) (*dto.ServiceListResponse, error) {

	// 1. Validasi Batas Halaman
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// 2. Konversi Page ke Offset untuk SQL
	limit := perPage
	offset := (page - 1) * perPage

	// 3. Panggil Repository
	services, totalItems, err := s.serviceRepo.FindAll(ctx, limit, offset, search, status, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// 4. Mapping dari Model ke DTO Summary
	var serviceResponses []dto.ServiceSummaryResponse
	for _, svc := range services {
		serviceResponses = append(serviceResponses, dto.ServiceSummaryResponse{
			ID:            svc.ID,
			Code:          svc.Code,
			ServiceName:   svc.ServiceName,
			Unit:          svc.Unit,
			Price:         svc.Price,
			DurationHours: svc.DurationHours,
			IsActive:      svc.IsActive,
			Category: &dto.NestedCategoryResponse{
				ID:                  svc.CategoryID,
				CategoryName:        svc.CategoryName,
				CategoryDescription: svc.CategoryDescription,
			},
		})
	}

	// 5. Cegah nilai "null" di JSON jika database kosong
	if serviceResponses == nil {
		serviceResponses = []dto.ServiceSummaryResponse{}
	}

	// 6. Hitung Total Halaman (Menggunakan integer math standar Kapten)
	var totalPages int
	if perPage > 0 {
		totalPages = int((totalItems + int64(perPage) - 1) / int64(perPage))
	}

	// 7. Kembalikan Response Akhir
	return &dto.ServiceListResponse{
		Data: serviceResponses,
		Meta: response.MetaData{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

// GetServiceDetail retrieves detailed service information by ID.
func (s *serviceService) GetServiceDetail(ctx context.Context, id int64) (*dto.ServiceDetailResponse, error) {

	// 1. Ambil Data dari Repository
	svc, err := s.serviceRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. Kembalikan Response menggunakan Helper
	return s.mapToDetailResponse(svc), nil
}

// ModifyService updates service profile with validation logic.
func (s *serviceService) ModifyService(ctx context.Context, targetID int64, req dto.UpdateServiceRequest) (*dto.ServiceDetailResponse, error) {

	// 1. Ambil Data Layanan yang Lama
	existingService, err := s.serviceRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	// 2. Validasi & Update Kode (Jika dikirim user)
	if req.Code != nil && *req.Code != existingService.Code {
		duplicateCheck, _ := s.serviceRepo.FindByCode(ctx, *req.Code)
		if duplicateCheck != nil && duplicateCheck.ID != targetID {
			return nil, response.ErrDuplicate
		}
		existingService.Code = *req.Code
	}

	// 3. Validasi & Update Nama (Jika dikirim user)
	if req.ServiceName != nil && *req.ServiceName != existingService.ServiceName {
		duplicateCheck, _ := s.serviceRepo.FindByName(ctx, *req.ServiceName)
		if duplicateCheck != nil && duplicateCheck.ID != targetID {
			return nil, response.ErrDuplicate
		}
		existingService.ServiceName = *req.ServiceName
	}

	// 4. Update Fields Lainnya (Partial Update)
	if req.CategoryID != nil {
		existingService.CategoryID = *req.CategoryID
	}
	if req.Unit != nil {
		existingService.Unit = *req.Unit
	}
	if req.Price != nil {
		existingService.Price = *req.Price
	}
	if req.DurationHours != nil {
		existingService.DurationHours = *req.DurationHours
	}
	if req.IsActive != nil {
		existingService.IsActive = *req.IsActive
	}

	// 5. Update Waktu (Timestamp)
	now := time.Now()
	existingService.UpdatedAt = &now

	// 6. Buat Model murni untuk dikirim ke Repository (tanpa data Category JOIN)
	updateModel := &models.Service{
		ID:            existingService.ID,
		CategoryID:    existingService.CategoryID,
		Code:          existingService.Code,
		ServiceName:   existingService.ServiceName,
		Unit:          existingService.Unit,
		Price:         existingService.Price,
		DurationHours: existingService.DurationHours,
		IsActive:      existingService.IsActive,
		UpdatedAt:     existingService.UpdatedAt,
	}

	// 7. Simpan Perubahan ke Database
	err = s.serviceRepo.UpdateService(ctx, updateModel)
	if err != nil {
		return nil, err
	}

	// 8. Panggil ulang GetServiceDetail agar data JOIN Kategori ikut diperbarui
	return s.GetServiceDetail(ctx, existingService.ID)
}

// DeactivateService handles soft deletion of a service.
func (s *serviceService) DeactivateService(ctx context.Context, targetID int64) error {

	// 1. Cek apakah layanan tersebut ada
	_, err := s.serviceRepo.FindByID(ctx, targetID)
	if err != nil {
		return err
	}

	// 2. Eksekusi Soft Delete
	err = s.serviceRepo.DeleteService(ctx, targetID)
	if err != nil {
		return err
	}

	return nil
}

// --- HELPER FUNCTION ---

func (s *serviceService) mapToDetailResponse(svc *models.ServiceWithCategory) *dto.ServiceDetailResponse {

	// 1. Format CreatedAt menjadi string bersih
	createdAtStr := svc.CreatedAt.Format("2006-01-02 15:04:05")

	// 2. Format UpdatedAt menjadi pointer string (jika tidak nil)
	var updatedAtPtr *string
	if svc.UpdatedAt != nil {
		formatted := svc.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtPtr = &formatted
	}

	return &dto.ServiceDetailResponse{
		ID:            svc.ID,
		Code:          svc.Code,
		ServiceName:   svc.ServiceName,
		Unit:          svc.Unit,
		Price:         svc.Price,
		DurationHours: svc.DurationHours,
		IsActive:      svc.IsActive,
		CreatedAt:     createdAtStr, // Gunakan variabel string yang sudah diformat
		UpdatedAt:     updatedAtPtr, // Gunakan pointer string yang sudah diformat
		Category: &dto.NestedCategoryResponse{
			ID:                  svc.CategoryID,
			CategoryName:        svc.CategoryName,
			CategoryDescription: svc.CategoryDescription,
		},
	}
}
