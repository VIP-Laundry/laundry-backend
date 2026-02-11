package dto

// LoginRequest defines the payload for user login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse defines the payload returned upon successful login.
type LoginResponse struct {
	Token TokenResponse   `json:"token"`
	User  AuthUserSummary `json:"user"`
}

// TokenResponse contains the access token, refresh token, and expiration details.
type TokenResponse struct {
	TokenType    string `json:"token_type"`    // e.g., "Bearer"
	AccessToken  string `json:"access_token"`  // JWT
	RefreshToken string `json:"refresh_token"` // UUID
	ExpiresIn    int    `json:"expires_in"`    // Seconds
}

// AuthUserSummary holds minimal user info for the login response.
// Used by Frontend to store basic session info (e.g. in LocalStorage).
type AuthUserSummary struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"` // Added: Frontend usually needs this for display
	Username string `json:"username"`
	Role     string `json:"role"`
}

// RefreshTokenRequest defines the payload for requesting a new access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse returns a new access token.
// Note: We currently reuse the existing Refresh Token (until it expires).
type RefreshTokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// AuthMeResponse defines the user profile structure for the /auth/me endpoint.
type AuthMeResponse struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	Username    string `json:"username"`
	Email       string `json:"email"` // Added: Usually profile needs email
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"` // Added: Complete profile info
	IsActive    bool   `json:"is_active"`    // Changed: int -> bool (Consistency with User Module)
	CreatedAt   string `json:"created_at"`
}
