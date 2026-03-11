package routes

import (
	"laundry-backend/internal/config"
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes mengatur semua endpoint untuk modul kategori layanan.
func SetupCategoryRoutes(router *gin.RouterGroup, categoryHandler *handlers.CategoryHandler, authRepo repositories.AuthRepository, cfg *config.Config) {

	// Grouping URL: /api/v1/categories
	categories := router.Group("/categories")

	// Global Auth Middleware: Semua request ke /categories/* wajib bawa JWT valid
	categories.Use(middleware.AuthMiddleware(authRepo, cfg))

	// --- RESTRICTED ENDPOINTS (Hanya Owner & Cashier) ---
	categories.GET("", middleware.RoleMiddleware("owner", "cashier"), categoryHandler.HandleGetCategoryList)
	categories.GET("/:id", middleware.RoleMiddleware("owner", "cashier"), categoryHandler.HandleGetCategoryDetail)

	// --- STRICT RESTRICTED ENDPOINTS (Hanya Owner) ---
	categories.POST("", middleware.RoleMiddleware("owner"), categoryHandler.HandleCreateCategory)
	categories.PUT("/:id", middleware.RoleMiddleware("owner"), categoryHandler.HandleUpdateCategory)
	categories.DELETE("/:id", middleware.RoleMiddleware("owner"), categoryHandler.HandleDeleteCategory)
}
