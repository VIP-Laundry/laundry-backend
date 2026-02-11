package middleware

import (
	"laundry-backend/internal/config"
	"laundry-backend/internal/repositories"
	"laundry-backend/pkg/response"
	"laundry-backend/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware adalah penjaga gerbang untuk memvalidasi Access Token (JWT).
func AuthMiddleware(authRepo repositories.AuthRepository, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")

		// 2. Cek apakah header kosong atau format salah
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			// REFACTOR: Pakai Helper
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Unauthorized: Missing or invalid token format", "Authorization header is required")
			return
		}

		// 3. Potong string "Bearer " untuk mendapatkan token murni
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Ambil Secret Key dari Config
		secretKey := []byte(cfg.JWT.Secret)

		// 4. Validasi token menggunakan jwt_utils yang sudah kita buat
		claims, err := utils.ValidateAccessToken(tokenString, secretKey)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrInvalidToken, "Unauthorized: Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// 5. [BARU] Cek apakah JTI token ini ada di daftar Blacklist database
		isBlacklisted, err := authRepo.IsBlacklisted(c.Request.Context(), claims.ID)
		if err != nil {
			// REFACTOR: Pakai Helper
			response.ErrorResponse(c, http.StatusInternalServerError, response.ErrInternalServer, "Internal server error during auth check", nil)
			return
		}

		if isBlacklisted {
			// REFACTOR: Pakai Helper (Bisa pakai ErrUnauthorized atau bikin ErrTokenBlacklisted di codes.go)
			// Disini kita pakai ErrUnauthorized agar aman
			response.ErrorResponse(c, http.StatusUnauthorized, response.ErrUnauthorized, "Unauthorized: Token has been logged out", "Please login again")
			return
		}

		// 6. [BARU] Simpan data penting ke Context agar bisa dipakai di Handler Logout
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("jti", claims.ID)             // Penting untuk Logout
		c.Set("exp", claims.ExpiresAt.Time) // Penting untuk Logout

		// 7. Lanjutkan ke proses berikutnya
		c.Next()
	}
}
