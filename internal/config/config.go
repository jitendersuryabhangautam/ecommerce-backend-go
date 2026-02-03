package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTExpiry time.Duration

	AllowedOrigins []string

	StockReservationTTL time.Duration
}

func LoadConfig() *Config {
	// Load .env file
	godotenv.Load()

	// Parse JWT expiry
	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	// Parse stock reservation TTL
	stockTTLMinutes, _ := strconv.Atoi(getEnv("STOCK_RESERVATION_TTL_MINUTES", "10"))

	// Parse allowed origins (comma-separated)
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	origins := []string{}
	for _, o := range strings.Split(allowedOrigins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000"}
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "ecommerce_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		JWTExpiry: time.Duration(jwtExpiryHours) * time.Hour,

		AllowedOrigins: origins,

		StockReservationTTL: time.Duration(stockTTLMinutes) * time.Minute,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
