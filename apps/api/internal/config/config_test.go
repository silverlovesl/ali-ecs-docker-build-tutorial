package config_test

import (
	"testing"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
)

func TestLoadNormalizesBasePath(t *testing.T) {
	cases := map[string]string{
		"":               "",
		"/":              "",
		"tutorial":       "/tutorial",
		"/tutorial":      "/tutorial",
		"/tutorial/":     "/tutorial",
		"  /tutorial/  ": "/tutorial",
		"/deep/tutorial": "/deep/tutorial",
	}

	for input, want := range cases {
		t.Setenv("WEB_BASE_PATH", input)

		if got := config.Load().BasePath; got != want {
			t.Errorf("WEB_BASE_PATH=%q: expected %q, got %q", input, want, got)
		}
	}
}
