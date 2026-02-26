package middlewares

import (
	"laundry-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware memastikan user memiliki role yang sesuai untuk mengakses endpoint.
// Param `allowedRoles` bisa lebih dari satu, misal: "owner", "cashier".
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		// 1. Ambil role user dari Context (yang sudah dipasang oleh AuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists {
			// [FIX 1 & 2] Pakai CodeUnauthorized dan panggil c.Abort()
			response.ErrorResponse(c, http.StatusUnauthorized, response.CodeUnauthorized, "Unauthorized", "User role not found in context")
			c.Abort()
			return
		}

		// [FIX 3] Safe Type Assertion: Mencegah Panic jika tipe data bukan string
		roleStr, ok := userRole.(string)
		if !ok {
			response.ErrorResponse(c, http.StatusInternalServerError, response.CodeInternalServer, "Internal Server Error", "Invalid role type in context")
			c.Abort()
			return
		}

		// 2. Cek apakah role user ada di dalam daftar allowedRoles
		roleAllowed := false
		for _, role := range allowedRoles {
			if role == roleStr {
				roleAllowed = true
				break
			}
		}

		// 3. Jika tidak cocok, tolak akses (403 Forbidden)
		if !roleAllowed {
			// [FIX 1 & 2] Pakai CodeForbidden dan panggil c.Abort()
			response.ErrorResponse(c, http.StatusForbidden, response.CodeForbidden, "Forbidden Access", "You do not have permission to access this resource")
			c.Abort()
			return
		}

		// 4. Lanjut ke Handler utama jika lulus
		c.Next()
	}
}
