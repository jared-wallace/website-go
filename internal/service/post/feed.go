package post

import (
	"context"

	"github.com/jared-wallace/website-go/internal/model"
)

// ListForFeed returns the most recent N published posts for RSS feed generation.
// Unlike ListPublished (which returns PostSummary projections), this returns
// the raw model.Post slice so the RSS handler has access to RenderedHTML.
func (s *Service) ListForFeed(ctx context.Context, limit int) ([]model.Post, error) {
	return s.repo.ListPublished(ctx, limit, 0)
}

// ListSlugsForSitemap returns all published posts. The caller only needs Slug
// and CreatedAt, but we reuse repo.ListPublished to avoid a new interface method.
// At blog scale (hundreds of posts) fetching all rows is fine.
func (s *Service) ListSlugsForSitemap(ctx context.Context) ([]model.Post, error) {
	return s.repo.ListPublished(ctx, 10000, 0)
}
