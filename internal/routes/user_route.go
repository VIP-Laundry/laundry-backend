package routes

import (
	"laundry-backend/internal/config" // [TAMBAH INI] Untuk membaca secret key JWT
	"laundry-backend/internal/handlers"
	middleware "laundry-backend/internal/middlewares"
	"laundry-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures the HTTP routes for user management.
// It applies authentication and role-based access control (RBAC) middleware.
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, authRepo repositories.AuthRepository, cfg *config.Config) { // [FIX] Tambahkan cfg di sini

	// Create a route group for "/users"
	users := router.Group("/users")

	// Apply Global Authentication Middleware.
	// This ensures that all requests to /users/* must have a valid JWT token.
	// It injects "user_id" and "role" into the Gin Context.
	// [FIX] Masukkan 'cfg' ke AuthMiddleware
	users.Use(middleware.AuthMiddleware(authRepo, cfg))

	// 1. Create User (POST /api/v1/users)
	// Access: Restricted to 'owner' only.
	users.POST("", middleware.RoleMiddleware("owner"), userHandler.CreateUser)

	// 2. Get User Directory (GET /api/v1/users)
	// Access: Restricted to 'owner' only.
	users.GET("", middleware.RoleMiddleware("owner"), userHandler.GetListUsers)

	// 3. Get User Detail (GET /api/v1/users/:id)
	// Access: Restricted to 'owner' only.
	users.GET("/:id", middleware.RoleMiddleware("owner"), userHandler.GetDetailUser)

	// 4. Update User Profile (PUT /api/v1/users/:id)
	// Access: Owner, Cashier, Staff, Courier.
	users.PUT("/:id", middleware.RoleMiddleware("owner", "cashier", "staff", "courier"), userHandler.UpdateUser)

	// 5. Delete/Deactivate User (DELETE /api/v1/users/:id)
	// Access: Restricted to 'owner' only.
	users.DELETE("/:id", middleware.RoleMiddleware("owner"), userHandler.DeleteUser)
}
