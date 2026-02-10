package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ==========================================
// 1. STRUCT DEFINITIONS (RESPONSE DTO)
// ==========================================

// BaseResponse defines the standard JSON structure for all API responses.
// This structure ensures consistency across success and error responses.
type BaseResponse struct {
	Success bool        `json:"success"`        // Success indicates whether the request was processed successfully.
	Message string      `json:"message"`        // Message provides a human-readable description of the operation result.
	Data    interface{} `json:"data"`           // Data contains the requested payload. It can be a single object, an array, or null.
	Meta    interface{} `json:"meta,omitempty"` // Meta contains pagination information. It is omitted if empty.
}

// MetaData defines the structure for pagination details.
// It is used within the 'Meta' field of BaseResponse.
type MetaData struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// ErrorResponseData defines the structure for detailed error reporting.
// It is typically put inside the 'Data' field when Success is false.
type ErrorResponseData struct {
	ErrorCode string      `json:"error_code"` // ErrorCode is a standardized machine-readable code (e.g., "VALIDATION_ERROR").
	Errors    interface{} `json:"errors"`     // Errors contains specific validation messages or debugging info.
}

// ==========================================
// 2. HELPER FUNCTION RESPONSE
// ==========================================

// SuccessOK sends a standard success response with HTTP 200 status code.
// Use this for successful GET (Detail), PUT, and DELETE operations.
func SuccessOK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessCreated sends a standard success response with HTTP 201 status code.
// Use this specifically for successful POST (Creation) operations.
func SuccessCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a standard success response containing pagination metadata.
// Use this for list endpoints (GET) that support pagination.
func SuccessMeta(c *gin.Context, message string, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// ErrorResponse sends a standardized error response.
// This function acts as a general error handler for various HTTP status codes (400, 401, 404, 500, etc.).
//
// Parameters:
//   - code: HTTP Status Code (e.g., http.StatusBadRequest).
//   - errorCode: A standardized string code (e.g., "VALIDATION_ERROR", "DUPLICATE_DATA").
//   - message: A human-readable error message in English.
//   - errors: Detailed error info (can be nil, string, or struct).
func ErrorResponse(c *gin.Context, code int, errorCode string, message string, errors interface{}) {
	c.JSON(code, BaseResponse{
		Success: false,
		Message: message,
		Data: ErrorResponseData{
			ErrorCode: errorCode,
			Errors:    errors,
		},
	})
	// Abort ensures that no further handlers are executed in the Gin middleware chain.
	c.Abort()
}
