package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"laundry-backend/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Asia%%2FJakarta",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka driver database: %v", err)
	}

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	lifetime := time.Duration(cfg.DB.MaxLifetimeMin) * time.Minute
	db.SetConnMaxLifetime(lifetime)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database tidak merespon (cek apakah MySQL Laragon sudah nyala?): %v", err)
	}

	log.Println("✅ Koneksi Database Berhasil Terhubung!")
	return db, nil
}
