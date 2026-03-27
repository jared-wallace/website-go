package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jared-wallace/website-go/internal/model"
)

// ErrSlugExists is returned by Create when the slug violates the unique constraint.
var ErrSlugExists = errors.New("post with that slug already exists")

// pgUniqueViolation is the Postgres error code for unique constraint violations.
const pgUniqueViolation = "23505"

// FindByID returns the post with the given ID, or ErrNotFound if absent.
func (r *postgresRepository) FindByID(ctx context.Context, id int64) (*model.Post, error) {
	const q = `
		SELECT id, title, slug, body, rendered_html, tags, published, created_at, updated_at, deleted_at
		FROM posts
		WHERE id = $1`

	var p model.Post
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Body,
		&p.RenderedHTML,
		&p.Tags,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &p, nil
}

// ListAll returns all posts in reverse chronological order, including soft-deleted ones.
func (r *postgresRepository) ListAll(ctx context.Context) ([]model.Post, error) {
	const q = `
		SELECT id, title, slug, body, rendered_html, tags, published, created_at, updated_at, deleted_at
		FROM posts
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("ListAll query: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Slug,
			&p.Body,
			&p.RenderedHTML,
			&p.Tags,
			&p.Published,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("ListAll scan: %w", err)
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListAll rows: %w", err)
	}
	return posts, nil
}

// Create inserts a new post and returns the record populated with ID and timestamps.
// Returns ErrSlugExists if the slug column unique constraint is violated.
func (r *postgresRepository) Create(ctx context.Context, p model.Post) (*model.Post, error) {
	const q = `
		INSERT INTO posts (title, slug, body, rendered_html, tags, published)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, q,
		p.Title,
		p.Slug,
		p.Body,
		p.RenderedHTML,
		p.Tags,
		p.Published,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return nil, ErrSlugExists
		}
		return nil, fmt.Errorf("Create: %w", err)
	}
	return &p, nil
}

// Update modifies a non-deleted post identified by p.ID.
func (r *postgresRepository) Update(ctx context.Context, p model.Post) error {
	const q = `
		UPDATE posts
		SET title = $1, slug = $2, body = $3, rendered_html = $4, tags = $5, updated_at = now()
		WHERE id = $6 AND deleted_at IS NULL`

	_, err := r.pool.Exec(ctx, q,
		p.Title,
		p.Slug,
		p.Body,
		p.RenderedHTML,
		p.Tags,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("Update: %w", err)
	}
	return nil
}

// SoftDelete sets deleted_at on the post, hiding it from public views.
func (r *postgresRepository) SoftDelete(ctx context.Context, id int64) error {
	const q = `UPDATE posts SET deleted_at = now(), updated_at = now() WHERE id = $1`
	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("SoftDelete: %w", err)
	}
	return nil
}

// Restore clears deleted_at and sets published=false so the post re-enters draft state.
func (r *postgresRepository) Restore(ctx context.Context, id int64) error {
	const q = `UPDATE posts SET deleted_at = NULL, published = false, updated_at = now() WHERE id = $1`
	if _, err := r.pool.Exec(ctx, q, id); err != nil {
		return fmt.Errorf("Restore: %w", err)
	}
	return nil
}

// SetPublished toggles the published flag on a non-deleted post.
func (r *postgresRepository) SetPublished(ctx context.Context, id int64, published bool) error {
	const q = `UPDATE posts SET published = $1, updated_at = now() WHERE id = $2 AND deleted_at IS NULL`
	if _, err := r.pool.Exec(ctx, q, published, id); err != nil {
		return fmt.Errorf("SetPublished: %w", err)
	}
	return nil
}
