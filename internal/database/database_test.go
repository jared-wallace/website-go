package database

import (
	"context"
	"os"
	"testing"
)

// integration tests require a live Postgres instance.
// Run with: DATABASE_URL=postgres://... go test ./internal/database/... -v
// In CI, the Postgres service container provides this automatically.

func TestConnect_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Postgres — skipping in short mode")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	ctx := context.Background()
	pool, err := Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer pool.Close()

	if pool == nil {
		t.Fatal("expected non-nil pool")
	}
}

func TestConnect_BadDSN(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Postgres — skipping in short mode")
	}

	ctx := context.Background()
	pool, err := Connect(ctx, "postgres://invalid:invalid@localhost:9999/nonexistent?connect_timeout=2")
	if err == nil {
		pool.Close()
		t.Fatal("expected error for bad DSN, got nil")
	}
}

func TestRunMigrations(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Postgres — skipping in short mode")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	ctx := context.Background()
	pool, err := Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer pool.Close()

	if migErr := RunMigrations(ctx, pool); migErr != nil {
		t.Fatalf("RunMigrations: %v", migErr)
	}

	// Verify the posts table was created.
	var exists bool
	err = pool.QueryRow(ctx,
		"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'posts')",
	).Scan(&exists)
	if err != nil {
		t.Fatalf("querying information_schema: %v", err)
	}
	if !exists {
		t.Fatal("expected posts table to exist after migrations")
	}
}
