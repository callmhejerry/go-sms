package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/callmhejerry/sms/internal/identity"
	"github.com/callmhejerry/sms/internal/shared/auth"
	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/tenant"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

type Handlers struct {
	Tenant   *tenant.Handler
	Identity *identity.Handler
}

func New(port string, pool *pgxpool.Pool, logger *slog.Logger, handlers Handlers, jwtManger *auth.JWTManager, roleChecker middleware.RoleChecker) *Server {
	mux := http.NewServeMux()

	// --------------------------
	// Public routes (no auth)
	// --------------------------
	healthHanler := NewHealthHandler(pool)
	mux.HandleFunc("GET /healthz", healthHanler.Healthz)
	mux.HandleFunc("GET /readyz", healthHanler.Readyz)

	mux.HandleFunc("POST /api/v1/tenants", handlers.Tenant.CreateTenant)
	mux.HandleFunc("POST /api/v1/login", handlers.Identity.Login)

	// --------------------------
	// Protected routes
	// --------------------------
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/v1/tenants", handlers.Tenant.ListTenants)
	protectedMux.HandleFunc("GET /api/v1/tenants/{id}", handlers.Tenant.GetTenant)

	// IDENTITY ROUTE

	adminOnly := middleware.RequireRole(roleChecker, "admin")

	protectedMux.Handle("POST /api/v1/users", adminOnly(http.HandlerFunc(handlers.Identity.CreateUser)))
	protectedMux.HandleFunc("GET /api/v1/tenants/{tenant_id}/users/{id}", handlers.Identity.GetUser)

	// MIDDLEWARE CHAIN
	protectedHandler := middleware.AuthMiddleware(jwtManger)(protectedMux)
	mux.Handle("/", protectedHandler)

	var handler http.Handler = mux

	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(logger)(handler)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	server := &Server{
		httpServer: httpServer,
		logger:     logger,
	}
	return server
}

func (server *Server) Start() error {
	server.logger.Info("Starting HTTP SERVER", slog.String("addr", server.httpServer.Addr))
	return server.httpServer.ListenAndServe()
}

func (server *Server) ShutDown(ctx context.Context) error {
	server.logger.Info("Shutting down HTTP server")
	return server.httpServer.Shutdown(ctx)
}
