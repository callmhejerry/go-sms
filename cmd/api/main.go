package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/callmhejerry/sms/internal/shared/config"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/logger"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()

	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.AppEnv)

	log.Info("Starting application", slog.String("env", cfg.AppEnv), slog.String("port", cfg.AppPort))

	fmt.Println("Database DSN:", cfg.DSN())

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg)

	if err != nil {
		log.Error("failed to connect to database", slog.String("error", err.Error()))
		panic(err)
	}
	defer pool.Close()
	log.Info("Connected to the database successfully")
	fmt.Println("Database connection OK")
}
