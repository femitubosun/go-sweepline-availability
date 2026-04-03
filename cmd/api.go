package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/femitubosun/go-sweepline-availability/internal/config"
	"github.com/femitubosun/go-sweepline-availability/internal/location"
	"github.com/femitubosun/go-sweepline-availability/internal/redis"
)

type app struct {
	config   *config.Config
	logger   *slog.Logger
	services services
}

type services struct {
	location location.Service
}

func (a *app) mount() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	a.registerStaticRoute(mux)
	a.registerLocationRoutes(mux)

	return mux
}

func (a *app) run(ctx context.Context, h http.Handler) error {
	addr := fmt.Sprintf(":%d", a.config.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		WriteTimeout: a.config.WriteTimeout,
		ReadTimeout:  a.config.ReadTimeout,
		IdleTimeout:  a.config.IdleTimeout,
	}

	errChan := make(chan error, 1)

	go func() {
		a.logger.Info("🌱 server started", "addr", addr)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		a.logger.Info("🍂 shutting down server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)

}

func NewApp(cfg *config.Config, cache *redis.Cache, logger *slog.Logger) *app {
	return &app{config: cfg, logger: logger, services: services{
		location: location.NewService(),
	}}
}
