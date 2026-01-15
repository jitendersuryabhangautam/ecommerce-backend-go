package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"ecommerce-backend/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := DB.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL database")

	// Run migrations
	if err := runMigrations(ctx); err != nil {
		log.Printf("⚠️ Warning: Could not run migrations: %v", err)
	}

	return DB, nil
}

func runMigrations(ctx context.Context) error {
	// In production, use a proper migration tool like golang-migrate
	// For now, we'll just log that migrations should be run manually
	log.Println("📋 Please run the SQL migrations from the migrations/ directory")
	return nil
}

func GetDB() *pgxpool.Pool {
	return DB
}

// BeginTx starts a new transaction
func BeginTx(ctx context.Context) (pgx.Tx, error) {
	return DB.Begin(ctx)
}

// HealthCheck checks database health
func HealthCheck(ctx context.Context) error {
	return DB.Ping(ctx)
}
