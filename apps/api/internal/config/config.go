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
	// WebDir is a directory holding a built single-page app. When set (and it
	// contains index.html) the server also serves the frontend, which is how the
	// production image ships as a single process. Empty disables static serving
	// so that `vite dev` keeps owning the frontend during development.
	WebDir string
	// BasePath is an optional URL prefix the app is mounted under, e.g. "/tutorial".
	// It lets one image serve correctly both directly at the root and behind a
	// reverse proxy that forwards the prefix unchanged. Empty means root.
	BasePath string
}

// Load reads configuration from environment variables and applies safe defaults
// for local development.
func Load() Config {
	return Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		AllowOrigins: splitAndTrim(getEnv("ALLOW_ORIGINS", "http://localhost:3000")),
		WebDir:       getEnv("WEB_DIR", ""),
		BasePath:     normalizeBasePath(getEnv("WEB_BASE_PATH", "")),
	}
}

// normalizeBasePath turns a configured prefix into a canonical "/tutorial" form:
// leading slash, no trailing slash, and "" for the root.
func normalizeBasePath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return ""
	}

	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}

	return strings.TrimRight(trimmed, "/")
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
