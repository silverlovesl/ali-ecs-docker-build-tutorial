// Package middleware contains reusable Gin middleware.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS allows the configured browser origins to call the API. When allowOrigins
// is empty every origin is accepted, which is convenient for local development.
func CORS(allowOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowOrigins))
	for _, origin := range allowOrigins {
		allowed[origin] = struct{}{}
	}
	allowAll := len(allowed) == 0

	return func(c *gin.Context) {
		if origin := c.GetHeader("Origin"); origin != "" {
			if _, ok := allowed[origin]; ok || allowAll {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				c.Header("Vary", "Origin")
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
