package response

// Standardized Error Codes for the Application.
// Frontend can use these strings to determine logic flow (e.g., show login modal if UNAUTHORIZED).
const (

	// ============================================
	// 1. GENERAL ERRORS (Masalah Umum)
	// ============================================
	ErrValidation     = "VALIDATION_ERROR"      // ErrValidation: Input data is missing or invalid format. (400 Bad Request)
	ErrUnauthorized   = "UNAUTHORIZED_ACCESS"   // ErrUnauthorized: Missing or invalid authentication token. (401 Unauthorized)
	ErrForbidden      = "FORBIDDEN_ACCESS"      // ErrForbidden: Authenticated but insufficient permissions (e.g., User accessing Owner menu). (403 Forbidden)
	ErrNotFound       = "RESOURCE_NOT_FOUND"    // ErrNotFound: The requested resource ID was not found in database. (404 Not Found)
	ErrDuplicate      = "DUPLICATE_DATA"        // ErrDuplicate: Unique constraint violation (e.g., Email/Username already exists). (409 Conflict)
	ErrRateLimit      = "RATE_LIMIT_EXCEEDED"   // ErrRateLimit: Too many requests sent in short period. (429 Too Many Requests)
	ErrInternalServer = "INTERNAL_SERVER_ERROR" // ErrInternal: Unexpected server or database failure. (500 Internal Server Error)

	// ============================================
	// 2. AUTH SPECIFIC ERRORS (Masalah Login)
	// ============================================
	ErrInvalidCredentials = "INVALID_CREDENTIALS"
	ErrAccountInactive    = "ACCOUNT_INACTIVE"
	ErrTokenExpired       = "TOKEN_EXPIRED"
	ErrInvalidToken       = "INVALID_TOKEN"
	ErrUserNotFound       = "USER_NOT_FOUND"
)
