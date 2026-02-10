package main

import (
	"fmt"
	"laundry-backend/internal/config"
	"laundry-backend/internal/db"
	"log"
)

func main() {
	// 1. Load Konfigurasi (Baca .env)
	fmt.Println("1. Memuat konfigurasi...")
	cfg := config.LoadConfig()

	// 2. Tes Koneksi Database
	fmt.Println("2. Mencoba menghubungi database...")
	database, err := db.ConnectDB(cfg)
	if err != nil {
		// Jika gagal, program mati di sini dan kasih tau errornya
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}

	// Jangan lupa tutup koneksi kalau program selesai
	defer database.Close()

	// 3. Jika kode sampai sini, berarti BERHASIL
	fmt.Println("------------------------------------------------")
	fmt.Println("✅ SUKSES! Database MySQL terhubung dengan lancar.")
	fmt.Printf("   - Host: %s\n", cfg.DB.Host)
	fmt.Printf("   - Database: %s\n", cfg.DB.Name)
	fmt.Println("------------------------------------------------")

	// Program selesai, tidak perlu menjalankan Server Gin dulu
}
