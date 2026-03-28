package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jared-wallace/website-go/internal/model"
)

// ListPublished returns published, non-deleted posts in reverse chronological order.
// Results are bounded by limit and offset for pagination.
func (r *postgresRepository) ListPublished(ctx context.Context, limit, offset int) ([]model.Post, error) {
	const q = `
		SELECT id, title, slug, body, rendered_html, tags, published, created_at, updated_at
		FROM posts
		WHERE published = true AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("ListPublished query: %w", err)
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
		); err != nil {
			return nil, fmt.Errorf("ListPublished scan: %w", err)
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListPublished rows: %w", err)
	}
	return posts, nil
}

// CountPublished returns the number of published, non-deleted posts.
func (r *postgresRepository) CountPublished(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(*) FROM posts WHERE published = true AND deleted_at IS NULL`
	var count int
	if err := r.pool.QueryRow(ctx, q).Scan(&count); err != nil {
		return 0, fmt.Errorf("CountPublished: %w", err)
	}
	return count, nil
}

// FindBySlug returns the post with the given slug, or ErrNotFound when absent.
func (r *postgresRepository) FindBySlug(ctx context.Context, slug string) (*model.Post, error) {
	const q = `
		SELECT id, title, slug, body, rendered_html, tags, published, created_at, updated_at
		FROM posts
		WHERE slug = $1 AND deleted_at IS NULL
		LIMIT 1`

	var p model.Post
	err := r.pool.QueryRow(ctx, q, slug).Scan(
		&p.ID,
		&p.Title,
		&p.Slug,
		&p.Body,
		&p.RenderedHTML,
		&p.Tags,
		&p.Published,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("FindBySlug: %w", err)
	}
	return &p, nil
}
