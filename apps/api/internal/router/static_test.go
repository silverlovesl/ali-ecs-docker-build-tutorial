package router_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/router"
)

const indexMarker = `data-marker="spa-shell"`

// newDist writes a minimal frontend build and returns its directory.
func newDist(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	files := map[string]string{
		"index.html":      `<!doctype html><html><body><div id="root" ` + indexMarker + `></div></body></html>`,
		"assets/app.js":   `console.log("bundled")`,
		"favicon.svg":     `<svg xmlns="http://www.w3.org/2000/svg"/>`,
		"nested/deep.txt": "deep",
	}

	for name, body := range files {
		full := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
		if err := os.WriteFile(full, []byte(body), 0o600); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	return dir
}

func newEngine(t *testing.T, webDir string) http.Handler {
	t.Helper()

	return newEngineAt(t, webDir, "")
}

// newEngineAt builds a handler that is also mounted under basePath.
func newEngineAt(t *testing.T, webDir, basePath string) http.Handler {
	t.Helper()

	gin.SetMode(gin.TestMode)

	return router.New(config.Config{
		Env:          "test",
		Port:         "0",
		AllowOrigins: []string{"http://localhost:3000"},
		WebDir:       webDir,
		BasePath:     basePath,
	})
}

func get(t *testing.T, engine http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

	return recorder
}

func TestStaticServesIndexAtRoot(t *testing.T) {
	recorder := get(t, newEngine(t, newDist(t)), "/")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), indexMarker) {
		t.Fatalf("expected the SPA shell, got %q", recorder.Body.String())
	}
	if got := recorder.Header().Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("expected index.html to be revalidated, got Cache-Control %q", got)
	}
}

func TestStaticFallsBackToIndexForClientRoutes(t *testing.T) {
	engine := newEngine(t, newDist(t))

	for _, path := range []string{"/about", "/settings/profile", "/deeply/nested/route"} {
		recorder := get(t, engine, path)

		if recorder.Code != http.StatusOK {
			t.Fatalf("%s: expected status %d, got %d", path, http.StatusOK, recorder.Code)
		}
		if !strings.Contains(recorder.Body.String(), indexMarker) {
			t.Fatalf("%s: expected the SPA shell to be served", path)
		}
	}
}

func TestStaticServesRealFiles(t *testing.T) {
	engine := newEngine(t, newDist(t))

	asset := get(t, engine, "/assets/app.js")
	if asset.Code != http.StatusOK || !strings.Contains(asset.Body.String(), "bundled") {
		t.Fatalf("expected the bundled asset, got %d %q", asset.Code, asset.Body.String())
	}
	if got := asset.Header().Get("Cache-Control"); !strings.Contains(got, "immutable") {
		t.Fatalf("expected hashed assets to be immutable, got Cache-Control %q", got)
	}

	for _, path := range []string{"/favicon.svg", "/nested/deep.txt"} {
		recorder := get(t, engine, path)
		if recorder.Code != http.StatusOK {
			t.Fatalf("%s: expected status %d, got %d", path, http.StatusOK, recorder.Code)
		}
	}
}

func TestStaticMissingHashedAssetIs404(t *testing.T) {
	engine := newEngine(t, newDist(t))

	recorder := get(t, engine, "/assets/stale-hash.js")

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
	if body := recorder.Body.String(); strings.Contains(body, indexMarker) {
		t.Fatalf("a missing hashed asset must not be answered with the SPA shell, got %q", body)
	}
}

func TestStaticKeepsAPIErrorsAsJSON(t *testing.T) {
	recorder := get(t, newEngine(t, newDist(t)), "/tutorial-api/does-not-exist")

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
	if body := recorder.Body.String(); strings.Contains(body, indexMarker) {
		t.Fatalf("the SPA fallback must not mask API 404s, got %q", body)
	}
	if contentType := recorder.Header().Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected a JSON error, got Content-Type %q", contentType)
	}
}

func TestStaticDoesNotShadowAPI(t *testing.T) {
	recorder := get(t, newEngine(t, newDist(t)), "/tutorial-api/health")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), indexMarker) {
		t.Fatalf("the API route must win over static serving, got %q", recorder.Body.String())
	}
}

func TestStaticServesUnderBasePath(t *testing.T) {
	engine := newEngineAt(t, newDist(t), "/tutorial")

	// The SPA shell answers under the prefix, both with and without the slash.
	for _, path := range []string{"/tutorial", "/tutorial/"} {
		recorder := get(t, engine, path)
		if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), indexMarker) {
			t.Fatalf("%s: expected the SPA shell, got %d %q", path, recorder.Code, recorder.Body.String())
		}
	}

	// Hashed assets keep their immutable cache header when fetched via the prefix.
	asset := get(t, engine, "/tutorial/assets/app.js")
	if asset.Code != http.StatusOK || !strings.Contains(asset.Body.String(), "bundled") {
		t.Fatalf("expected the bundled asset under /tutorial/, got %d %q", asset.Code, asset.Body.String())
	}
	if got := asset.Header().Get("Cache-Control"); !strings.Contains(got, "immutable") {
		t.Fatalf("expected hashed assets to be immutable, got Cache-Control %q", got)
	}

	// API routes are reachable under the prefix too.
	if code := get(t, engine, "/tutorial/tutorial-api/health").Code; code != http.StatusOK {
		t.Fatalf("expected the API under /tutorial/, got %d", code)
	}

	// Unknown prefixed API paths stay JSON instead of falling back to the shell.
	missing := get(t, engine, "/tutorial/tutorial-api/nope")
	if missing.Code != http.StatusNotFound || strings.Contains(missing.Body.String(), indexMarker) {
		t.Fatalf("expected a JSON 404 under /tutorial/, got %d %q", missing.Code, missing.Body.String())
	}

	// A missing hashed asset under the prefix is a real 404, not the shell.
	stale := get(t, engine, "/tutorial/assets/stale.js")
	if stale.Code != http.StatusNotFound || strings.Contains(stale.Body.String(), indexMarker) {
		t.Fatalf("expected a 404 for a stale prefixed asset, got %d %q", stale.Code, stale.Body.String())
	}
}

func TestStaticBasePathOnlyMatchesPathBoundary(t *testing.T) {
	dir := newDist(t)
	// A sibling directory whose name starts with the mount prefix.
	if err := os.MkdirAll(filepath.Join(dir, "tutorials"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tutorials", "list.txt"), []byte("list"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	engine := newEngineAt(t, dir, "/tutorial")

	recorder := get(t, engine, "/tutorials/list.txt")
	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), "list") {
		t.Fatalf("expected /tutorials/list.txt to be served, got %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestStaticDisabledWithoutIndexHTML(t *testing.T) {
	// Local development: Vite owns the frontend, so unmatched paths stay 404.
	engine := newEngine(t, t.TempDir())

	for _, path := range []string{"/", "/about"} {
		if code := get(t, engine, path).Code; code != http.StatusNotFound {
			t.Fatalf("%s: expected status %d, got %d", path, http.StatusNotFound, code)
		}
	}
}

func TestStaticIgnoresMissingWebDir(t *testing.T) {
	engine := newEngine(t, filepath.Join(t.TempDir(), "does-not-exist"))

	if code := get(t, engine, "/").Code; code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, code)
	}
}
