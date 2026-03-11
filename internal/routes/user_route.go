package routes

import (
	"laundry-backend/internal/config"
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes mengatur semua endpoint untuk manajemen pengguna (User Directory).
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authRepo repositories.AuthRepository, cfg *config.Config) {

	// Grouping URL: /api/v1/users
	users := router.Group("/users")

	// Global Auth Middleware: Semua request ke /users/* wajib bawa JWT valid
	users.Use(middleware.AuthMiddleware(authRepo, cfg))

	// --- RESTRICTED ENDPOINTS (Role-Based Access) ---
	// Create & Read (Hanya Owner)
	users.POST("", middleware.RoleMiddleware("owner"), userHandler.CreateUser)
	users.GET("", middleware.RoleMiddleware("owner"), userHandler.GetListUsers)
	users.GET("/:id", middleware.RoleMiddleware("owner"), userHandler.GetDetailUser)

	// Update (Bisa diakses staff lain untuk update profil sendiri/operasional)
	users.PUT("/:id", middleware.RoleMiddleware("owner", "cashier", "staff", "courier"), userHandler.UpdateUser)

	// Delete/Deactivate (Hanya Owner)
	users.DELETE("/:id", middleware.RoleMiddleware("owner"), userHandler.DeleteUser)
}
