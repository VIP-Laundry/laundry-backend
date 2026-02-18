package response

import "errors"

// ============================================
// 1. STRING CODES (Untuk Frontend/JSON)
// ============================================
// Kita gunakan prefix 'Msg' agar tidak tertukar dengan objek error.
// Ini menjaga agar Modul Auth atau bagian lain yang butuh string tetap aman.
const (
	CodeValidation     = "VALIDATION_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED_ACCESS"
	CodeForbidden      = "FORBIDDEN_ACCESS"
	CodeNotFound       = "RESOURCE_NOT_FOUND"
	CodeDuplicate      = "DUPLICATE_DATA"
	CodeRateLimit      = "RATE_LIMIT_EXCEEDED"
	CodeInternalServer = "INTERNAL_SERVER_ERROR"

	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeAccountInactive    = "ACCOUNT_INACTIVE"
	CodeTokenExpired       = "TOKEN_EXPIRED"
	CodeInvalidToken       = "INVALID_TOKEN"
	CodeUserNotFound       = "USER_NOT_FOUND"
)

// ============================================
// 2. SENTINEL ERRORS (Untuk Internal Go/Repository/Service)
// ============================================
// Ini yang akan kita kembalikan di fungsi return.
// Mengikuti standar E (Error Handling) & O (Optimal Go Practice).
var (
	ErrValidation     = errors.New(CodeValidation)
	ErrUnauthorized   = errors.New(CodeUnauthorized)
	ErrForbidden      = errors.New(CodeForbidden)
	ErrNotFound       = errors.New(CodeNotFound)
	ErrDuplicate      = errors.New(CodeDuplicate)
	ErrRateLimit      = errors.New(CodeRateLimit)
	ErrInternalServer = errors.New(CodeInternalServer)

	ErrInvalidCredentials = errors.New(CodeInvalidCredentials)
	ErrAccountInactive    = errors.New(CodeAccountInactive)
	ErrTokenExpired       = errors.New(CodeTokenExpired)
	ErrInvalidToken       = errors.New(CodeInvalidToken)
	ErrUserNotFound       = errors.New(CodeUserNotFound)
)
