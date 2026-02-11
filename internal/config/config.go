package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APP  AppConfig
	DB   DBConfig
	JWT  JWTConfig
	CORS CORSConfig
	LOG  LOGConfig
}

type AppConfig struct {
	Name  string
	Env   string
	Port  string
	Debug bool
}

type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	MaxIdleConns   int
	MaxOpenConns   int
	MaxLifetimeMin int
}

type JWTConfig struct {
	Secret            string
	ExpiryMin         int
	RefreshExpiryHour int
}

type CORSConfig struct {
	AllowedOrigins string
}

type LOGConfig struct {
	Level string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultValue
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	return &Config{
		APP: AppConfig{
			Name:  getEnv("APP_NAME", "VIP Laundry Backend"),
			Env:   getEnv("APP_ENV", "development"),
			Port:  getEnv("APP_PORT", "8080"),
			Debug: getEnvBool("APP_DEBUG", true),
		},
		DB: DBConfig{
			Host:           getEnv("DB_HOST", "127.0.0.1"),
			Port:           getEnv("DB_PORT", "3306"),
			User:           getEnv("DB_USER", "root"),
			Password:       getEnv("DB_PASSWORD", ""),
			Name:           getEnv("DB_NAME", "viplaundry"),
			MaxOpenConns:   getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetimeMin: getEnvAsInt("DB_MAX_LIFETIME_MINUTES", 5),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", ""),
			ExpiryMin:         getEnvAsInt("JWT_EXPIRY_MINUTE", 15),
			RefreshExpiryHour: getEnvAsInt("JWT_REFRESH_EXPIRY_HOUR", 24),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
		LOG: LOGConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}
