package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamNilotpal/k8s-demo/internal/config"
	"github.com/iamNilotpal/k8s-demo/internal/server"
	"github.com/iamNilotpal/k8s-demo/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables in development.
	if os.Getenv(config.EnvLookupKey) == config.EnvDevelopment {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("error loading envs : %+v", err)
			os.Exit(1)
		}
	}

	// Load configuration.
	cfg := config.Load()
	fmt.Printf("Configuration loaded successfully : %+v \n", cfg)

	// Initialize structured logging with observability context.
	log := logger.NewWithTracing(cfg.ServiceName)
	defer func() {
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
	}()

	log.Infow("Starting k8s-demo platform with observability...")

	// Run the application with proper error handling.
	if err := run(log, cfg); err != nil {
		log.Errorw("startup error", "error", err)
		if err := log.Sync(); err != nil {
			log.Infow("sync error", "error", err)
		}
		os.Exit(1)
	}
}

func run(log *logger.Logger, cfg *config.AppConfig) error {
	// Create application context for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize and start the server with observability.
	srv, err := server.New(ctx, cfg, log)
	if err != nil {
		return err
	}

	// Setup graceful shutdown.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine.
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.Start()
	}()

	// Wait for shutdown signal or server error.
	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		log.Infow("shutting down server", "signal", sig)
		return srv.Shutdown(ctx)
	}
}
