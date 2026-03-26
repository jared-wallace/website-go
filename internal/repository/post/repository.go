// Package post provides the repository layer for blog post persistence.
// It defines the Repository interface and a Postgres-backed implementation.
package post

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jared-wallace/website-go/internal/model"
)

// ErrNotFound is returned by FindBySlug when no matching post exists.
var ErrNotFound = errors.New("post not found")

// Repository is the data-access contract for blog posts. Callers depend on
// this interface rather than the concrete Postgres type, enabling mock
// implementations in unit tests.
type Repository interface {
	// ListPublished returns published, non-deleted posts ordered by created_at DESC.
	ListPublished(ctx context.Context, limit, offset int) ([]model.Post, error)
	// CountPublished returns the total number of published, non-deleted posts.
	CountPublished(ctx context.Context) (int, error)
	// FindBySlug returns the post matching slug, or ErrNotFound if none exists.
	FindBySlug(ctx context.Context, slug string) (*model.Post, error)
}

// postgresRepository implements Repository against a pgxpool connection pool.
type postgresRepository struct {
	pool *pgxpool.Pool
}

// New constructs a Repository backed by the supplied Postgres connection pool.
func New(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}
