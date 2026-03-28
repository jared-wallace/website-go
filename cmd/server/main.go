package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/database"
	"github.com/jared-wallace/website-go/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()
	logger.Info("config loaded", "env", cfg.AppEnv, "port", cfg.Port)

	ctx, stop := server.GracefulShutdown()
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connected")

	if err := database.RunMigrations(ctx, pool); err != nil {
		logger.Error("migrations failed", "error", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")

	// TODO(Phase 2): wire HTTP handler and register routes.
	// srv := server.New(cfg.Port, mux)

	logger.Info("server starting", "port", cfg.Port)

	// Block until signal received (Phase 2 will replace with srv.ListenAndServe).
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = shutdownCtx // Phase 2 will pass this to srv.Shutdown(shutdownCtx)

	logger.Info("server stopped")
}
