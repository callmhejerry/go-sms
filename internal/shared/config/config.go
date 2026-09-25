package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost     string
	DBUser     string
	DBName     string
	DBPort     string
	DBPassword string
	DBSSLMode  string

	JWTSecret          string
	JWTExpirationHours int
}

func Load() (*Config, error) {
	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8081"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBUser:             getEnv("DB_USER", "school"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBPassword:         getEnv("DB_PASSWORD", "schoolsecret"),
		DBName:             getEnv("DB_NAME", "school_sms"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production-super-secret-key"),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 72),
	}
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("Invalid config: %w", err)
	}
	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if cfg.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	// Stronger rules in production
	if cfg.AppEnv == "production" {
		if len(cfg.JWTSecret) < 32 {
			return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
		}
		if cfg.DBSSLMode == "disable" {
			return fmt.Errorf("DB_SSLMODE should not be 'disable' in production")
		}
		if cfg.JWTSecret == "change-me-in-production-super-secret-key" ||
			cfg.JWTSecret == "change-me-in-production-super-secret-key-32chars" {
			return fmt.Errorf("JWT_SECRET must be changed from the default value in production")
		}
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
