package handlers

import (
	"laundry-backend/internal/dto"
	"laundry-backend/internal/services"
	"laundry-backend/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user authentication.
// @Summary User Login
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	// 1. Validate Input
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Invalid input format", err.Error())
		return
	}

	// 2. Call Service
	res, err := h.authService.AuthenticateUser(c.Request.Context(), req)
	if err != nil {
		// Mapping Error string dari Service ke response.ErrorResponse
		switch err.Error() {
		case response.ErrInvalidCredentials:
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrInvalidCredentials, "Invalid username or password", nil)
		case response.ErrAccountInactive:
			response.ErrorResponse(c, http.StatusForbidden, response.ErrAccountInactive, "Your account is inactive", nil)
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected error occurred", err.Error())
		}
		return
	}

	// 3. Success Response
	response.SuccessOK(c, "Login successfully", res)
}

// RefreshToken handles access token renewal.
// @Summary Refresh Access Token
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {

	var req dto.RefreshTokenRequest

	// 1. Validate Input
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Refresh token is required", err.Error())
		return
	}

	// 2. Call Service
	res, err := h.authService.RenewUserSession(c.Request.Context(), req)
	if err != nil {
		switch err.Error() {
		case response.ErrInvalidToken:
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrInvalidToken, "Invalid or revoked token", nil)
		case response.ErrTokenExpired:
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrTokenExpired, "Session expired, please login again", nil)
		case response.ErrUserNotFound:
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUserNotFound, "User account not found", nil)
		default:
			response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "An unexpected error occurred", err.Error())
		}
		return
	}

	// 3. Success Response
	response.SuccessOK(c, "Access token refreshed successfully", res)
}

// Logout handles user session termination.
// @Summary User Logout
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {

	var req dto.RefreshTokenRequest

	// 1. Bind Refresh Token from Body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorResponse(c, http.StatusBadRequest, response.ErrValidation, "Refresh token is required for logout", err.Error())
		return
	}

	// 2. Extract JTI and Exp from Context (Set by Middleware)
	jti, existsJti := c.Get("jti")
	exp, existsExp := c.Get("exp")
	userID, existsUserID := c.Get("user_id")

	if !existsJti || !existsExp || !existsUserID {
		response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Invalid session context", nil)
		return
	}

	// 3. Call Service
	err := h.authService.RevokeUserSession(c.Request.Context(), req.RefreshToken, jti.(string), exp.(time.Time), userID.(int64))

	if err != nil {
		// Handle Security Check (Jika mencoba logout token orang lain)
		if err.Error() == response.ErrUnauthorized {
			response.ErrorResponse(c, http.StatusForbidden, response.ErrForbidden, "You cannot logout another user's session", nil)
			return
		}

		// Handle error umum lainnya
		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to process logout", err.Error())
		return
	}

	// 4. Success Response
	response.SuccessOK(c, "Logout successfully", gin.H{"status": "access_terminated"})
}

// GetMe retrieves the authenticated user's profile.
// @Summary Get Current User Profile
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {

	// 1. Extract UserID from Context (Set by Middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "User context missing", nil)
		return
	}

	// 2. Call Service
	res, err := h.authService.GetAccountProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		if err.Error() == response.ErrAccountInactive {
			response.ErrorResponse(c, http.StatusForbidden, response.ErrAccountInactive, "Your account is inactive", nil)
			return
		}

		response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Failed to retrieve profile", err.Error())
		return
	}

	// 3. Success Response
	response.SuccessOK(c, "User profile retrieved successfully", res)
}
