package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/router"
)

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := router.New(config.Config{
		Env:          "test",
		Port:         "0",
		AllowOrigins: []string{"http://localhost:3000"},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if payload["status"] != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", payload["status"])
	}
}

func TestCORSPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := router.New(config.Config{
		Env:          "test",
		Port:         "0",
		AllowOrigins: []string{"http://localhost:3000"},
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected allow origin header, got %q", got)
	}
}
