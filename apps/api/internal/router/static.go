package router

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// apiPrefix is the route group reserved for JSON endpoints.
const apiPrefix = "/tutorial-api"

// assetsPrefix holds the content-hashed bundle emitted by Vite.
const assetsPrefix = "/assets/"

// registerStaticFiles serves a built single-page app from dir. Real files are
// returned as-is; any other path falls back to index.html so that client-side
// routes survive a hard refresh.
//
// Two things are deliberately not swallowed by the SPA fallback:
//
//   - unknown paths under /tutorial-api, which keep answering JSON so that a
//     typo in an API call is not masked by a 200 HTML response
//   - non-GET/HEAD requests
//
// dir must contain index.html. When it does not, static serving stays disabled
// so that local development (Vite dev server in front of this API) is
// unaffected.
func registerStaticFiles(engine *gin.Engine, dir string) {
	if dir == "" {
		return
	}

	indexPath := filepath.Join(dir, "index.html")
	if info, err := os.Stat(indexPath); err != nil || info.IsDir() {
		log.Printf("WEB_DIR=%s has no index.html; serving API routes only", dir)
		return
	}

	engine.NoRoute(func(c *gin.Context) {
		requestPath := path.Clean("/" + c.Request.URL.Path)

		if requestPath == apiPrefix || strings.HasPrefix(requestPath, apiPrefix+"/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		filePath := filepath.Join(dir, filepath.FromSlash(strings.TrimPrefix(requestPath, "/")))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			setCacheControl(c, requestPath)
			c.File(filePath)
			return
		}

		// Vite hashes the names of everything under /assets, so a miss there is a
		// genuine 404 (the previous nginx used `try_files $uri =404`) rather than
		// a client-side route.
		if strings.HasPrefix(requestPath, assetsPrefix) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Client-side route: hand back the shell and let the router resolve it.
		c.Header("Cache-Control", "no-cache")
		c.File(indexPath)
	})
}

// setCacheControl reproduces what the previous nginx sidecar did: Vite writes
// content-hashed filenames under /assets, so those can be cached forever, while
// everything else (index.html above all) must be revalidated on each load.
func setCacheControl(c *gin.Context, requestPath string) {
	if strings.HasPrefix(requestPath, assetsPrefix) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		return
	}

	c.Header("Cache-Control", "no-cache")
}
