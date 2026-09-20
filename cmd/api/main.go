package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/callmhejerry/sms/internal/academic"
	"github.com/callmhejerry/sms/internal/admission"
	"github.com/callmhejerry/sms/internal/identity"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/config"
	"github.com/callmhejerry/sms/internal/shared/database"
	"github.com/callmhejerry/sms/internal/shared/logger"
	"github.com/callmhejerry/sms/internal/shared/server"
	"github.com/callmhejerry/sms/internal/shared/store"
	"github.com/callmhejerry/sms/internal/shared/validation"
	"github.com/callmhejerry/sms/internal/tenant"
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

	queries := store.New(pool)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpirationHours)

	validator := validation.NewValdiator()

	// Services
	tenantService := tenant.NewService(pool, queries)
	identityService := identity.NewService(queries, jwtManager)
	academicService := academic.NewService(queries, pool)
	admissionService := admission.NewService(queries, academicService)

	// Handlers
	tenantHandler := tenant.NewHandler(tenantService, log)
	identityHandler := identity.NewHandler(identityService, log)
	academicHandler := academic.NewHandler(academicService, log)
	admissionHandler := admission.NewHandler(admissionService, log, validator)

	handlers := server.Handlers{
		Tenant:    tenantHandler,
		Identity:  identityHandler,
		Academic:  academicHandler,
		Admission: admissionHandler,
	}

	server := server.New(cfg.AppPort, pool, log, handlers, jwtManager, identityService)
	//graceful shutdown
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("Sever error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// wait for interupt signal
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown signal received")

	shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.ShutDown(shutdownctx); err != nil {
		log.Error("server shutdown failed", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
}
