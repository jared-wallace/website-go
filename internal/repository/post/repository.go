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

	// FindByID returns the post with the given ID, or ErrNotFound if absent.
	FindByID(ctx context.Context, id int64) (*model.Post, error)
	// ListAll returns all posts (including unpublished and soft-deleted) ordered by created_at DESC.
	ListAll(ctx context.Context) ([]model.Post, error)
	// Create inserts a new post and returns the fully populated record (with ID and timestamps).
	// Returns ErrSlugExists if the slug is already taken.
	Create(ctx context.Context, p model.Post) (*model.Post, error)
	// Update modifies an existing, non-deleted post identified by p.ID.
	Update(ctx context.Context, p model.Post) error
	// SoftDelete marks a post as deleted without removing it from the database.
	SoftDelete(ctx context.Context, id int64) error
	// Restore clears the deleted_at timestamp and sets published=false on a soft-deleted post.
	Restore(ctx context.Context, id int64) error
	// SetPublished toggles the published flag on a non-deleted post.
	SetPublished(ctx context.Context, id int64, published bool) error
}

// postgresRepository implements Repository against a pgxpool connection pool.
type postgresRepository struct {
	pool *pgxpool.Pool
}

// New constructs a Repository backed by the supplied Postgres connection pool.
func New(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}
