package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/callmhejerry/sms/internal/shared/middleware"
	"github.com/callmhejerry/sms/internal/tenant"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

type Handlers struct {
	TenantHandler *tenant.Handler
}

func New(port string, pool *pgxpool.Pool, logger *slog.Logger, handlers Handlers) *Server {
	mux := http.NewServeMux()

	healthHanler := NewHealthHandler(pool)

	mux.HandleFunc("GET /healthz", healthHanler.Healthz)
	mux.HandleFunc("GET /readyz", healthHanler.Readyz)

	// TENANT ROUTE
	mux.HandleFunc("POST /api/v1/tenants", handlers.TenantHandler.CreateTenant)
	mux.HandleFunc("GET /api/v1/tenants", handlers.TenantHandler.ListTenants)
	mux.HandleFunc("GET /api/v1/tenants/{id}", handlers.TenantHandler.GetTenant)

	// MIDDLEWARE CHAIN
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
