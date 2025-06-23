package config

import (
	"os"
	"strconv"
	"time"
)

func getEnvIntOrFallback(key string, fallback int) int {
	env, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	if parsed, err := strconv.Atoi(env); err == nil {
		return parsed
	}
	return fallback
}

func getEnvOrFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationOrFallback(key, fallback string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(fallback)
	return duration
}
