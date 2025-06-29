package handlers

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	health_handlers "github.com/iamNilotpal/k8s-demo/internal/handlers/health"
	"github.com/iamNilotpal/k8s-demo/internal/metrics"
	"github.com/iamNilotpal/k8s-demo/internal/middlewares"
	"github.com/iamNilotpal/k8s-demo/pkg/logger"
)

type Config struct {
	Service string
	Version string
	DB      *sql.DB
	Router  *chi.Mux
	Log     *logger.Logger
	Metrics *metrics.Metrics
}

func SetupRoutes(cfg *Config) {
	cfg.Router.Use(middleware.RequestID)
	cfg.Router.Use(middleware.RealIP)
	cfg.Router.Use(middleware.Logger)
	cfg.Router.Use(middleware.Recoverer)

	cfg.Router.Use(middlewares.MetricsMiddleware(cfg.Metrics))
	cfg.Router.Use(middlewares.TracingMiddleware(cfg.Service))

	healthHandlers := health_handlers.New(&health_handlers.Config{
		Service: cfg.Service,
		Version: cfg.Version,
		DB:      cfg.DB,
		Log:     cfg.Log,
		Metrics: cfg.Metrics,
	})

	cfg.Router.Handle("/metrics", promhttp.Handler())
	cfg.Router.Get("/health", healthHandlers.HealthCheck)
}
