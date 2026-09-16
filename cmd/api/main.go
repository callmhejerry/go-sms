package main

import (
	"fmt"
	"log/slog"

	"github.com/callmhejerry/sms/internal/shared/config"
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
}
