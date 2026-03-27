package post

import (
	"context"

	"github.com/jared-wallace/website-go/internal/model"
)

// Create renders body markdown, then inserts a new post via the repository.
// Returns ErrSlugExists (from the repository layer) if the slug is taken.
func (s *Service) Create(ctx context.Context, title, slug, body, tags string, publish bool) (*model.Post, error) {
	rendered := s.renderer.Render(body)
	p := model.Post{
		Title:        title,
		Slug:         slug,
		Body:         body,
		RenderedHTML: string(rendered),
		Tags:         tags,
		Published:    publish,
	}
	return s.repo.Create(ctx, p)
}

// Update re-renders body markdown, then persists changes to an existing post.
func (s *Service) Update(ctx context.Context, id int64, title, slug, body, tags string) error {
	rendered := s.renderer.Render(body)
	p := model.Post{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Body:         body,
		RenderedHTML: string(rendered),
		Tags:         tags,
	}
	return s.repo.Update(ctx, p)
}

// SoftDelete marks a post deleted without removing it from the database.
func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	return s.repo.SoftDelete(ctx, id)
}

// Restore clears the deleted state and resets the post to unpublished draft.
func (s *Service) Restore(ctx context.Context, id int64) error {
	return s.repo.Restore(ctx, id)
}

// Publish makes a post publicly visible.
func (s *Service) Publish(ctx context.Context, id int64) error {
	return s.repo.SetPublished(ctx, id, true)
}

// Unpublish hides a post from the public blog without deleting it.
func (s *Service) Unpublish(ctx context.Context, id int64) error {
	return s.repo.SetPublished(ctx, id, false)
}

// ListAll returns all posts including drafts and soft-deleted entries.
func (s *Service) ListAll(ctx context.Context) ([]model.Post, error) {
	return s.repo.ListAll(ctx)
}

// GetByID returns the post with the given ID, or ErrNotFound if absent.
func (s *Service) GetByID(ctx context.Context, id int64) (*model.Post, error) {
	return s.repo.FindByID(ctx, id)
}
