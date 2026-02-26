package handlers

import (
	"errors"
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"laundry-backend/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser handles POST /api/v1/users.
// Access: Owner only.
func (h *UserHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest

	// 1. Validate Input
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Input validation failed", err.Error())
		return
	}

	// 2. Call Service
	res, err := h.userService.RegisterUser(c.Request.Context(), req)

	if err != nil {
		if errors.Is(err, response.ErrDuplicate) {
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "User data already exists", "Username, email, or phone number is already taken")
			return
		}

		// Default Internal Error
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response
	response.SuccessCreated(c, "User account created successfully", res)
}

// GetListUsers handles GET /api/v1/users.
// Access: Owner only.
func (h *UserHandler) GetListUsers(c *gin.Context) {

	// 1. Parse Query Parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	search := c.Query("search")
	role := c.Query("role")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	// 2. Call Service
	res, err := h.userService.GetUsers(c.Request.Context(), page, perPage, search, role, status)
	if err != nil {
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response with Meta
	response.SuccessMeta(c, "Users retrieved successfully", res.Data, res.Meta)
}

// GetDetailUser handles GET /api/v1/users/:id.
// Access: Owner only.
func (h *UserHandler) GetDetailUser(c *gin.Context) {

	// 1. Parse ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID format", nil)
		return
	}

	// 2. Call Service
	res, err := h.userService.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		// [FIX] Gunakan errors.Is
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "User not found", nil)
			return
		}
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 3. Success Response
	response.SuccessOK(c, "User detail retrieved successfully", res)
}

// UpdateUser handles PUT /api/v1/users/:id.
// Access: Owner (All), Others (Self only).
func (h *UserHandler) UpdateUser(c *gin.Context) {

	// 1. Parse Target ID
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID format", nil)
		return
	}

	// 2. Extract Requester Info (From Auth Middleware)
	// Note: We assume AuthMiddleware sets "user_id" (float64/int64) and "role" (string).
	// We use Type Assertion to ensure safety.
	requesterIDRaw, okID := c.Get("user_id")
	requesterRole, okRole := c.Get("role")

	if !okID || !okRole {
		response.ErrorResponse(c, http.StatusUnauthorized, response.CodeUnauthorized, "Invalid authentication context", nil)
		return
	}

	requesterID, okAssert := requesterIDRaw.(int64)
	if !okAssert {
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Invalid user ID type (expected int64)", nil)
		return
	}

	// 3. Bind JSON Input
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Input validation failed", err.Error())
		return
	}

	// 4. Call Service with Requester Context
	res, err := h.userService.ModifyUserData(c.Request.Context(), targetID, req, requesterID, requesterRole.(string))
	if err != nil {
		// [FIX] Map Specific Errors menggunakan errors.Is
		if errors.Is(err, response.ErrForbidden) {
			response.ErrorResponse(c, http.StatusForbidden, response.CodeForbidden, "You do not have permission to modify this profile", nil)
			return
		}
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "User not found", nil)
			return
		}
		if errors.Is(err, response.ErrDuplicate) {
			response.ErrorResponse(c, http.StatusConflict, response.CodeDuplicate, "Data conflict detected (duplicate entry)", "Username, email, or phone number is already taken")
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 5. Success Response
	response.SuccessOK(c, "User updated successfully", res)
}

// DeleteUser handles DELETE /api/v1/users/:id.
// Access: Owner only (Anti-self delete logic in Service).
func (h *UserHandler) DeleteUser(c *gin.Context) {

	// 1. Parse Target ID
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.CodeValidation, "Invalid ID format", nil)
		return
	}

	// 2. Extract Requester ID (From Auth Middleware)
	requesterIDRaw, ok := c.Get("user_id")
	if !ok {
		response.ErrorResponse(c, http.StatusUnauthorized, response.CodeUnauthorized, "Invalid authentication context", nil)
		return
	}

	requesterID, okAssert := requesterIDRaw.(int64)
	if !okAssert {
		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Invalid user ID type (expected int64)", nil)
		return
	}

	// 3. Call Service
	if err := h.userService.DeactivateUserAccount(c.Request.Context(), targetID, requesterID); err != nil {
		// [FIX] Gunakan errors.Is
		if errors.Is(err, response.ErrForbidden) {
			response.ErrorResponse(c, http.StatusForbidden, response.CodeForbidden, "Action not permitted (cannot delete self)", nil)
			return
		}
		if errors.Is(err, response.ErrNotFound) {
			response.ErrorResponse(c, http.StatusNotFound, response.CodeNotFound, "User not found", nil)
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "An unexpected server error occurred", nil)
		return
	}

	// 4. Success Response
	// We return the ID of the deleted user as data.
	response.SuccessOK(c, "User account deactivated successfully", gin.H{"id": targetID})
}
