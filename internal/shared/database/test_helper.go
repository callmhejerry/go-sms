package database

import (
	"context"
	"os"
	"testing"

	"github.com/callmhejerry/sms/internal/shared/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	cfg := &config.Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "school"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBName:     getEnv("DB_NAME", "school_sms"),
		DBPassword: getEnv("DB_PASSWORD", "schoolsecret"),
		DBSSLMode:  "disable",
	}

	pool, err := NewPool(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Failed to create test pool: %v", err)
	}
	return pool
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
