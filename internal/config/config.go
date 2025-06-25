package config

import (
	"time"
)

type Web struct {
	APIHost         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DB struct {
	MaxIdleConns int
	MaxOpenConns int
	TLS          string
	Name         string
	User         string
	Host         string
	Password     string
	Scheme       string
}

type APIConfig struct {
	DB  *DB
	Web *Web
}

func Load() *APIConfig {
	return &APIConfig{
		DB: &DB{
			TLS:          getEnvOrFallback("DB_TLS", "disable"),
			User:         getEnvOrFallback("DB_USER", "postgres"),
			Name:         getEnvOrFallback("DB_NAME", "k8s-demo"),
			Host:         getEnvOrFallback("DB_HOST", "localhost"),
			Scheme:       getEnvOrFallback("DB_SCHEME", "postgres"),
			Password:     getEnvOrFallback("DB_PASSWORD", "password"),
			MaxIdleConns: getEnvIntOrFallback("DB_MAX_IDLE_CONN", 5),
			MaxOpenConns: getEnvIntOrFallback("DB_MAX_OPEN_CONN", 20),
		},
		Web: &Web{
			WriteTimeout:    getDurationOrFallback("SERVER_IDLE_TIMEOUT", ""),
			ReadTimeout:     getDurationOrFallback("SERVER_READ_TIMEOUT", ""),
			IdleTimeout:     getDurationOrFallback("SERVER_WRITE_TIMEOUT", ""),
			ShutdownTimeout: getDurationOrFallback("SERVER_SHUTDOWN_TIMEOUT", ""),
			APIHost:         getEnvOrFallback("SERVER_API_HOST", "localhost:8080"),
		},
	}
}
