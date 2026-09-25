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
		DBName:     "school_sms_test",
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

func cleanup(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(context.Background(), `
		TRUNCATE TABLE
			audit_logs,
			payment_allocations,
			payments,
			student_fees,
			fee_structures,
			fee_types,
			scores,
			results,
			teacher_assignments,
			class_subjects,
			subjects,
			assessment_types,
			inventory_issuances,
			stock_movements,
			inventory_items,
			inventory_categories,
			student_parents,
			parents,
			admissions,
			students,
			class_arms,
			classes,
			academic_sessions,
			user_roles,
			roles,
			users,
			tenants
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		t.Fatalf("failed to clean test database: %v", err)
	}
}
