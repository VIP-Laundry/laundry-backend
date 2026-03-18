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

type ServiceHandler struct {
	serviceService services.ServiceService
}

func NewServiceHandler(serviceService services.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

func (h *ServiceHandler) HandleCreateService(c *gin.Context) {

	var req dto.CreateServiceRequest

	// 1. Validasi Payload JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid request payload", err.Error())
		return
	}

	// 2. Eksekusi Service dengan membawa Context
	res, err := h.serviceService.CreateService(c.Request.Context(), req)
	if err != nil {
		// Gunakan Sentinel Error (ErrDuplicate) untuk pengecekan Internal Go
		if errors.Is(err, response.ErrDuplicate) {
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "Service code or name already exists", nil)
			return
		}

		// Log error untuk kebutuhan internal/debugging
		fmt.Printf("[ERROR] CreateService: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to create service", nil)
		return
	}

	// 3. Sukses
	response.SuccessCreated(c, "Service created successfully", res)

}

func (h *ServiceHandler) HandleGetServiceList(c *gin.Context) {

	// 1. Ambil nilai dari URL Query Parameters
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "10")
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	// ===============================================================
	// 🛡️ [TAMBAHAN VIP-7: ROLE-BASED FILTERING]
	// ===============================================================
	userRole := c.GetString("role")

	// Jika Kasir, PAKSA status menjadi "1" (Hanya Aktif).
	if userRole == "cashier" {
		status = "1"
	}
	// ===============================================================

	// 2. Konversi tipe data dengan "Safety Net" (Keamanan Ekstra)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Fallback aman
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 {
		perPage = 10 // Fallback aman
	}

	// 3. Panggil Koki (Service)
	res, err := h.serviceService.GetServiceList(c.Request.Context(), page, perPage, search, status, sortBy, sortOrder)
	if err != nil {
		fmt.Printf("[ERROR] GetServiceList: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to retrieve services", nil)
		return
	}

	// 4. Sukses dengan Meta (Pagination)
	response.SuccessMeta(c, "Services retrieved successfully", res.Data, res.Meta)
}

func (h *ServiceHandler) HandleGetServiceDetail(c *gin.Context) {

	// 1. Ambil ID dari URL Path (misal: /services/5)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Panggil Service
	res, err := h.serviceService.GetServiceDetail(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Service not found", nil)
			return
		}

		fmt.Printf("[ERROR] GetServiceDetail: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to retrieve service detail", nil)
		return
	}

	// 3. Sukses
	response.SuccessOK(c, "Service detail retrieved successfully", res)
}

func (h *ServiceHandler) HandleUpdateService(c *gin.Context) {

	// 1. Ambil ID dari URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Ambil Data JSON dari Body
	var req dto.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid request payload", err.Error())
		return
	}

	// 3. Panggil Koki (Service)
	res, err := h.serviceService.ModifyService(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Service not found", nil)
			return
		}

		if errors.Is(err, response.ErrDuplicate) {
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "Service code or name already taken", nil)
			return
		}

		fmt.Printf("[ERROR] ModifyService: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to update service", nil)
		return
	}

	// 4. Sukses
	response.SuccessOK(c, "Service updated successfully", res)
}

func (h *ServiceHandler) HandleDeleteService(c *gin.Context) {

	// 1. Ambil ID dari URL Path
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID Format", "ID must be a number")
		return
	}

	// 2. Panggil Koki (Service) untuk menonaktifkan layanan
	err = h.serviceService.DeactivateService(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "Service not found", nil)
			return
		}

		fmt.Printf("[ERROR] DeactivateService: %v\n", err)

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Failed to delete service", nil)
		return
	}

	// 3. Sukses
	response.SuccessOK(c, "Service deleted successfully", map[string]int64{"id": id})
}
