package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(port string, pool *pgxpool.Pool, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	healthHanler := NewHealthHandler(pool)

	mux.HandleFunc("GET /healthz", healthHanler.Healthz)
	mux.HandleFunc("GET /readyz", healthHanler.Readyz)

	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
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
