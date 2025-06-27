package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamNilotpal/k8s-demo/internal/config"
	"github.com/iamNilotpal/k8s-demo/internal/database"
	"github.com/iamNilotpal/k8s-demo/internal/handlers"
	"github.com/iamNilotpal/k8s-demo/internal/metrics"
	"github.com/iamNilotpal/k8s-demo/internal/tracing"
	"github.com/iamNilotpal/k8s-demo/pkg/logger"
)

// Server represents our HTTP server with all dependencies.
type Server struct {
	db         *sql.DB
	httpServer *http.Server
	logger     *logger.Logger
	metrics    *metrics.Metrics
	config     *config.AppConfig
	shutdown   func(context.Context) error
}

func New(ctx context.Context, cfg *config.AppConfig, log *logger.Logger) (*Server, error) {
	// Initialize tracing using values from the centralized configuration.
	shutdown, err := tracing.New(
		&tracing.TracingConfig{
			ServiceName:    cfg.ServiceName,
			ServiceVersion: cfg.ServiceVersion,
			Environment:    cfg.Environment,
			JaegerEndpoint: cfg.JaegerEndpoint,
		},
	)
	if err != nil {
		log.Warnw("Failed to initialize tracing", "error", err)
		shutdown = func(context.Context) error { return nil }
	} else {
		log.Infow("Tracing initialized successfully",
			"service", cfg.ServiceName,
			"version", cfg.ServiceVersion,
			"environment", cfg.Environment,
			"endpoint", cfg.JaegerEndpoint,
		)
	}

	// Initialize metrics collection for performance monitoring.
	appMetrics := metrics.New()
	log.Infow("Metrics initialized successfully")

	// Initialize database connection.
	db, err := database.Open(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify database connectivity during startup with tracing context.
	dbCtx, dbSpan := tracing.StartSpan(ctx, cfg.ServiceName, "startup_check")
	if err := database.StatusCheck(dbCtx, db, log); err != nil {
		dbSpan.End()
		return nil, fmt.Errorf("database status check failed: %w", err)
	}
	dbSpan.End()
	log.Infow("Database connection verified successfully")

	// Create HTTP server using centralized web configuration.
	router := chi.NewRouter()
	handlers.SetupRoutes(&handlers.Config{
		DB:      db,
		Log:     log,
		Router:  router,
		Metrics: appMetrics,
		Service: cfg.ServiceName,
		Version: cfg.ServiceVersion,
	})

	server := &http.Server{
		Handler:      router,
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	return &Server{
		httpServer: server,
		db:         db,
		logger:     log,
		config:     cfg,
		metrics:    appMetrics,
		shutdown:   shutdown,
	}, nil
}

// Start begins serving HTTP requests using configuration values.
func (s *Server) Start() error {
	s.logger.Infow("server starting with full observability",
		"address", s.httpServer.Addr,
		"service", s.config.ServiceName,
		"version", s.config.ServiceVersion,
		"environment", s.config.Environment,
		"read_timeout", s.config.Web.ReadTimeout,
		"write_timeout", s.config.Web.WriteTimeout,
		"idle_timeout", s.config.Web.IdleTimeout,
	)

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Shutdown gracefully terminates the server and all observability components.
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Infow("initiating graceful shutdown",
		"service", s.config.ServiceName,
		"shutdown_timeout", s.config.Web.ShutdownTimeout,
	)

	// Create timeout context using the configured shutdown timeout.
	shutdownCtx, cancel := context.WithTimeout(ctx, s.config.Web.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server first to stop accepting new requests.
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	// Close database connection to free up database resources.
	if err := s.db.Close(); err != nil {
		s.logger.Warnw("error closing database connection", "error", err)
	}

	// Shutdown tracing to flush any remaining spans to the tracing backend.
	if err := s.shutdown(shutdownCtx); err != nil {
		s.logger.Warnw("error shutting down tracing", "error", err)
	}

	s.logger.Infow("graceful shutdown completed successfully",
		"service", s.config.ServiceName,
	)

	return nil
}
