package config

import (
	"time"
)

const (
	EnvProduction  string = "PRODUCTION"
	EnvDevelopment string = "DEVELOPMENT"
	EnvLookupKey   string = "ENVIRONMENT"
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

type AppConfig struct {
	DB             *DB
	Web            *Web
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
}

func Load() *AppConfig {
	return &AppConfig{
		ServiceName:    getEnvOrFallback("SERVICE_NAME", "k8s-demo"),
		ServiceVersion: getEnvOrFallback("SERVICE_VERSION", "v0.1.0"),
		Environment:    getEnvOrFallback(EnvLookupKey, EnvDevelopment),
		JaegerEndpoint: getEnvOrFallback("JAEGER_ENDPOINT", "http://jaeger:4318/v1/traces"),
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
			APIHost:         getEnvOrFallback("SERVER_API_HOST", ":8080"),
			ReadTimeout:     getDurationOrFallback("SERVER_READ_TIMEOUT", "10s"),
			WriteTimeout:    getDurationOrFallback("SERVER_WRITE_TIMEOUT", "10s"),
			IdleTimeout:     getDurationOrFallback("SERVER_IDLE_TIMEOUT", "120s"),
			ShutdownTimeout: getDurationOrFallback("SERVER_SHUTDOWN_TIMEOUT", "20s"),
		},
	}
}
