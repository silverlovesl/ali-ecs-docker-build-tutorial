// Package config resolves runtime configuration from the environment.
package config

import (
	"os"
	"strings"
)

// Config holds the runtime configuration for the API server.
type Config struct {
	// Env is the application environment: development, staging or production.
	Env string
	// Port is the TCP port the HTTP server binds to.
	Port string
	// AllowOrigins is the list of browser origins permitted by CORS.
	AllowOrigins []string
}

// Load reads configuration from environment variables and applies safe defaults
// for local development.
func Load() Config {
	return Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		AllowOrigins: splitAndTrim(getEnv("ALLOW_ORIGINS", "http://localhost:3000")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))

	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			values = append(values, trimmed)
		}
	}

	return values
}
