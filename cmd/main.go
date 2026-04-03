package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/femitubosun/go-sweepline-availability/internal/config"
	"github.com/femitubosun/go-sweepline-availability/internal/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("error", "err", err)
		os.Exit(1)
	}

	cache, err := redis.NewCache(cfg)
	if err != nil {
		logger.Error("failed to connect to redis", "err", err)
	}

	logger.Info("redis connected 🚀")

	api := NewApp(cfg, cache, logger)

	if err := api.run(ctx, api.mount()); err != nil {
		logger.Error("error", "err", err)
		os.Exit(1)
	}
}
