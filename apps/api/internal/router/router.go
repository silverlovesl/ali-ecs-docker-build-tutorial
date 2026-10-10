// Package router wires middleware and routes into a Gin engine.
package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/handler"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/middleware"
)

// New builds the HTTP handler with the project middleware and routes registered.
//
// When the app is mounted under a URL prefix (config.BasePath), the returned
// handler rewrites prefixed requests to their root equivalent before Gin routes
// them. That keeps one image working both behind a proxy that strips the prefix
// (the ECS nginx does) and one that forwards it unchanged. Rewriting up front,
// rather than registering duplicate routes, avoids collisions in Gin's radix
// tree — a base path such as "/tutorial" would otherwise clash with the
// "/tutorial-api" route prefix.
func New(cfg config.Config) http.Handler {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowOrigins))

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	health := handler.NewHealthHandler(cfg)

	api := engine.Group("/tutorial-api")
	{
		api.GET("/health", health.Get)
		api.GET("/hello", health.Hello)
	}

	// In production the built frontend ships inside this same binary's image, so
	// the process serves both the SPA and the API on a single port.
	registerStaticFiles(engine, cfg.WebDir)

	if cfg.BasePath == "" {
		return engine
	}

	return stripBasePath(cfg.BasePath, engine)
}

// stripBasePath returns a handler that drops basePath from incoming requests
// before delegating to next. The boundary check keeps a sibling path such as
// /tutorials from being mistaken for a /tutorial mount.
func stripBasePath(basePath string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rewritten, ok := trimBasePath(r.URL.Path, basePath); ok {
			r.URL.Path = rewritten
			// Keep RawPath consistent so downstream URL handling does not use
			// the un-rewritten escaped form.
			r.URL.RawPath = ""
		}

		next.ServeHTTP(w, r)
	})
}

// trimBasePath removes basePath from requestPath. It reports whether the path
// was prefixed at all; "/tutorial" and "/tutorial/..." both collapse to root.
func trimBasePath(requestPath, basePath string) (string, bool) {
	if requestPath == basePath {
		return "/", true
	}

	if strings.HasPrefix(requestPath, basePath+"/") {
		return strings.TrimPrefix(requestPath, basePath), true
	}

	return requestPath, false
}
