package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/router"
)

func main() {
	cfg := config.Load()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router.New(cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("api listening on http://localhost:%s (env=%s)", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("api stopped")
}
