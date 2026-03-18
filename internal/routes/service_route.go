package routes

import (
	"laundry-backend/internal/config"
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupServiceRoutes mengatur semua endpoint untuk modul layanan (services).
func SetupServiceRoutes(router *gin.RouterGroup, serviceHandler *handlers.ServiceHandler, authRepo repositories.AuthRepository, cfg *config.Config) {

	// Grouping URL: /api/v1/services
	services := router.Group("/services")

	// Global Auth Middleware: Semua request ke /services/* wajib bawa JWT valid
	services.Use(middleware.AuthMiddleware(authRepo, cfg))

	// --- STRICT RESTRICTED ENDPOINTS (Hanya Owner) ---
	// Endpoint untuk Create, Update, dan Delete (Mengubah Data)
	services.POST("", middleware.RoleMiddleware("owner"), serviceHandler.HandleCreateService)
	services.PUT("/:id", middleware.RoleMiddleware("owner"), serviceHandler.HandleUpdateService)
	services.DELETE("/:id", middleware.RoleMiddleware("owner"), serviceHandler.HandleDeleteService)

	// --- RESTRICTED ENDPOINTS (Hanya Owner & Cashier) ---
	// Endpoint untuk GetList dan GetDetail (Membaca Data)
	services.GET("", middleware.RoleMiddleware("owner", "cashier"), serviceHandler.HandleGetServiceList)
	services.GET("/:id", middleware.RoleMiddleware("owner", "cashier"), serviceHandler.HandleGetServiceDetail)
}
