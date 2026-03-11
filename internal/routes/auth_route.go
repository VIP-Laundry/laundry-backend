package routes

import (
	"laundry-backend/internal/config"
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes mengatur semua endpoint untuk autentikasi.
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, authRepo repositories.AuthRepository, cfg *config.Config) {

	// Grouping URL: /api/v1/auth
	auth := router.Group("/auth")
	{
		// --- PUBLIC ENDPOINTS (Tanpa Login) ---
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// --- PRIVATE ENDPOINTS (Wajib Login) ---
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware(authRepo, cfg))
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/me", authHandler.GetMe)
		}
	}
}
