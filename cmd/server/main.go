package main

import (
	"fmt"

	"laundry-backend/internal/handlers"
	"laundry-backend/internal/repositories"
	"laundry-backend/internal/routes"
	"laundry-backend/internal/services"

	"laundry-backend/internal/config"
	"laundry-backend/internal/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// ==========================================
	// 1. Load Konfigurasi (Baca .env)
	// ==========================================
	fmt.Println("1. Memuat konfigurasi...")
	cfg := config.LoadConfig()

	// ==========================================
	// 2. Tes Koneksi Database
	// ==========================================
	fmt.Println("2. Mencoba menghubungi database...")
	dbConn, err := db.ConnectDB(cfg)
	if err != nil {
		// Jika gagal, program mati di sini dan kasih tau errornya
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}
	// Jangan lupa tutup koneksi kalau program selesai
	defer dbConn.Close()

	// ==========================================
	// 3. DEPENDENCY INJECTION (WIRING)
	// ==========================================

	// We inject dependencies from the bottom up: DB -> Repo -> Service -> Handler
	fmt.Println("3. Merakit komponen internal (Dependency Injection)...")

	// A. Repository Layer (Data Access)
	userRepo := repositories.NewUserRepository(dbConn)
	authRepo := repositories.NewAuthRepository(dbConn)

	// B. Service Layer (Business Logic)
	authService := services.NewAuthService(authRepo, userRepo, cfg)

	// C. Handler Layer (HTTP Transport)
	authHandler := handlers.NewAuthHandler(authService)

	// ==========================================
	// 4. SETUP SERVER & ROUTES
	// ==========================================
	fmt.Println("4. Menyiapkan jalur API (Routes)...")

	r := gin.Default()

	// Global Group
	v1 := r.Group("/api/v1")

	// Daftarkan Module Auth
	routes.SetupAuthRoutes(v1, authHandler, authRepo, cfg)

	// ==========================================
	// 5. START THE SERVER
	// ==========================================
	fmt.Println("------------------------------------------------")
	fmt.Printf("✅ VIP Laundry Backend menyala di port %s\n", cfg.APP.Port)
	fmt.Println("------------------------------------------------")

	if err := r.Run(":" + cfg.APP.Port); err != nil {
		log.Fatalf("❌ Gagal menyalakan server: %v", err)
	}
}
