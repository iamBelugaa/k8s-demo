package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/iamNilotpal/k8s-demo/internal/config"
	"github.com/iamNilotpal/k8s-demo/internal/database"
	"github.com/iamNilotpal/k8s-demo/pkg/logger"
	"github.com/iamNilotpal/k8s-demo/pkg/response"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	log := logger.New("k8s-demo")
	defer func() {
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
	}()

	if err := godotenv.Load(); err != nil {
		log.Fatalw("error loading envs", "error", err)
	}

	log.Infow("Starting k8s-demo platform...")

	if err := run(log); err != nil {
		log.Infow("startup error", "error", err)
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config.Load()
	log.Infow("Configuration loaded successfully")

	db, err := database.Open(cfg.DB)
	if err != nil {
		return err
	}
	log.Infow("Database connection opened successfully")

	if err := database.StatusCheck(context.Background(), db); err != nil {
		return err
	}
	log.Infow("Database status check completed successfully")

	router := chi.NewRouter()
	registerRoutes(log, router, db)

	server := http.Server{
		Handler:      router,
		Addr:         cfg.Web.APIHost,
		ReadTimeout:  cfg.Web.ReadTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	serverErrors := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infow("server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutting down server", "signal", sig)
		defer log.Infow("shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func registerRoutes(log *zap.SugaredLogger, r *chi.Mux, db *sql.DB) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infow("Request received", "request", r)

		if err := database.StatusCheck(context.Background(), db); err != nil {
			response.RespondError(
				w, http.StatusInternalServerError,
				http.StatusText(http.StatusInternalServerError), "StatusInternalServerError", nil,
			)
		}

		response.RespondSuccess(w, http.StatusOK, http.StatusText(http.StatusOK), nil)
	})
}
