package handlers

import (
	"errors"
	"fmt"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"laundry-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) HandleCreateCategory(c *gin.Context) {

	var req dto.CreateCategoryRequest

	// 1. Validasi Payload JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		// Gunakan String Code (CodeValidation) untuk Frontend
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid request payload", err.Error())
		return
	}

	// 2. Eksekusi Service dengan membawa Context
	result, err := h.categoryService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		// Gunakan Sentinel Error (ErrDuplicate) untuk pengecekan Internal Go
		if errors.Is(err, response.ErrDuplicate) {
			// Gunakan String Code (CodeDuplicate) untuk Frontend
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "Category name already exists", nil)
			return
		}

		// Log error untuk kebutuhan internal/debugging
		fmt.Printf("[ERROR] CreateCategory: %v\n", err)

		// [PERBAIKAN] Gunakan String Code (CodeInternalServer) untuk Frontend
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to create category", nil)
		return
	}

	// 3. Sukses
	response.SuccessCreated(c, "Category created successfully", result)
}

func (h *CategoryHandler) HandleGetCategoryList(c *gin.Context) {

	// 1. Ambil nilai dari URL Query Parameters
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	search := c.Query("search")
	status := c.Query("status") // <-- Ini asalnya bisa kosong (""), "1", atau "0"
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("order")

	// ===============================================================
	// 🛡️ [TAMBAHAN VIP-7: ROLE-BASED FILTERING]
	// ===============================================================
	// Ambil role user yang sedang login dari Gin Context
	// (Ini hasil titipan dari AuthMiddleware)
	userRole := c.GetString("role")

	// Jika Kasir, PAKSA status menjadi "1" (Hanya Aktif).
	// Meskipun Kasir iseng ngetik URL: ?status=0, kita timpa jadi 1.
	if userRole == "cashier" {
		status = "1"
	}
	// Jika Owner, biarkan variabel status apa adanya.
	// Owner bisa melihat semua (status=""), yang aktif saja (status="1"),
	// atau yang sudah dihapus saja (status="0").
	// ===============================================================

	// 2. Konversi tipe data dengan "Safety Net" (Keamanan Ekstra)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Fallback aman: jika user iseng masukin huruf atau angka minus, paksa jadi halaman 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = 10 // Fallback aman: paksa jadi 10 item per halaman
	}

	// 3. Panggil Koki (Service)
	// [PERBAIKAN 1] Gunakan h.categoryService
	result, err := h.categoryService.GetCategoryList(c.Request.Context(), page, perPage, search, status, sortBy, sortOrder)
	if err != nil {
		fmt.Printf("[ERROR] GetCategoryList: %v\n", err)

		// [PERBAIKAN 2] Gunakan String Code (CodeInternalServer) untuk balasan ke Frontend
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to retrieve categories", nil)
		return
	}

	// 4. Sukses dengan Meta (Pagination)
	response.SuccessMeta(c, "Categories retrieved successfully", result.Data, result.Meta)
}

func (h *CategoryHandler) HandleGetCategoryDetail(c *gin.Context) {

	// 1. Ambil ID dari URL Path (misal: /categories/5)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// [PERBAIKAN] Gunakan String Code (CodeValidation)
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Panggil Service
	// [PERBAIKAN] Gunakan h.categoryService
	result, err := h.categoryService.GetCategoryDetail(c.Request.Context(), id)
	if err != nil {
		// [PERBAIKAN] Gunakan errors.Is dan Sentinel Error (ErrNotFound)
		if errors.Is(err, response.ErrNotFound) {
			// [PERBAIKAN] Gunakan String Code (CodeNotFound)
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Category not found", nil)
			return
		}

		fmt.Printf("[ERROR] GetCategoryDetail: %v\n", err)

		// [PERBAIKAN] Gunakan String Code (CodeInternalServer)
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to retrieve category detail", nil)
		return
	}

	// 3. Sukses
	response.SuccessOK(c, "Category detail retrieved successfully", result)
}

func (h *CategoryHandler) HandleUpdateCategory(c *gin.Context) {

	// 1. Ambil ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// [PERBAIKAN] Gunakan String Code (CodeValidation)
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Ambil Data JSON dari Body
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// [PERBAIKAN] Gunakan String Code (CodeValidation)
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid request payload", err.Error())
		return
	}

	// 3. Panggil Koki (Service)
	// [PERBAIKAN] Gunakan h.categoryService
	result, err := h.categoryService.ModifyCategory(c.Request.Context(), id, req)
	if err != nil {
		// [PERBAIKAN] Gunakan errors.Is untuk cek Not Found
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Category not found", nil)
			return
		}

		// [PERBAIKAN] Gunakan errors.Is untuk cek Duplicate (Misal nama barunya bentrok)
		if errors.Is(err, response.ErrDuplicate) {
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "Category name already taken", nil)
			return
		}

		fmt.Printf("[ERROR] ModifyCategory: %v\n", err)

		// [PERBAIKAN] Gunakan String Code (CodeInternalServer)
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to update category", nil)
		return
	}

	// 4. Sukses
	response.SuccessOK(c, "Category updated successfully", result)
}

func (h *CategoryHandler) HandleDeleteCategory(c *gin.Context) {

	// 1. Ambil ID dari URL Path
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// [PERBAIKAN] Gunakan String Code (CodeValidation)
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Panggil Koki (Service) untuk menonaktifkan kategori
	// [PERBAIKAN] Gunakan h.categoryService
	err = h.categoryService.DeactivateCategory(c.Request.Context(), id)
	if err != nil {
		// [PERBAIKAN] Gunakan errors.Is untuk cek Not Found
		if errors.Is(err, response.ErrNotFound) {
			// [PERBAIKAN] Gunakan String Code (CodeNotFound)
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Category not found", nil)
			return
		}

		fmt.Printf("[ERROR] DeactivateCategory: %v\n", err)

		// [PERBAIKAN] Gunakan String Code (CodeInternalServer)
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to delete category", nil)
		return
	}

	// 3. Sukses
	// Catatan: Mengembalikan map berisi ID yang dihapus adalah praktik yang sangat bagus!
	response.SuccessOK(c, "Category deleted successfully", map[string]int64{"id": id})
}
