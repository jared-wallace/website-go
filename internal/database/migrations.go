package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	migrations "github.com/jared-wallace/website-go/db/migrations"
	"github.com/pressly/goose/v3"
)

// RunMigrations applies all pending goose migrations embedded in the
// db/migrations package. It uses the pgx stdlib adapter to satisfy goose's
// *sql.DB requirement without abandoning native pgx elsewhere.
func RunMigrations(_ context.Context, pool *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDBFromPool(pool)
	defer func() { _ = sqlDB.Close() }()

	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(sqlDB, "."); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}
