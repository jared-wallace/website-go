package main

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jared-wallace/website-go/internal/config"
	"github.com/jared-wallace/website-go/internal/database"
	bloghandler "github.com/jared-wallace/website-go/internal/handler/blog"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
	"github.com/jared-wallace/website-go/internal/server"
	"github.com/jared-wallace/website-go/web"
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

	// Build dependency graph: repo -> service -> handler
	repo := postrepo.New(pool)
	svc := postservice.New(repo)
	blog := bloghandler.New(svc)

	// Register routes
	mux := http.NewServeMux()
	staticFS, err := fs.Sub(web.Static, "static")
	if err != nil {
		logger.Error("static fs sub failed", "error", err)
		os.Exit(1)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(staticFS)))
	mux.HandleFunc("GET /{$}", blog.ListPosts)          // Home page (exact match)
	mux.HandleFunc("GET /posts", blog.ListPosts)         // /posts?page=N
	mux.HandleFunc("GET /posts/{slug}", blog.ShowPost)   // Single post
	mux.HandleFunc("GET /{path...}", blog.NotFound)      // Catch-all 404

	srv := server.New(cfg.Port, mux)
	logger.Info("server starting", "port", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
	logger.Info("server stopped")
}
