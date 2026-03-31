package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/femitubosun/go-sweepline-availability/internal/config"
)

type app struct {
	config *config.Config
	logger *slog.Logger
}

func (a *app) mount() http.Handler {

	return nil
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
		a.logger.Info("server started", "addr", addr)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		a.logger.Info("shuttdin down server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)

}

func NewApp(cfg *config.Config, logger *slog.Logger) *app {
	return &app{config: cfg, logger: logger}
}
