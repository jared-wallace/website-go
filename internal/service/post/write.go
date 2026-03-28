package post

import (
	"context"
	"errors"

	"github.com/jared-wallace/website-go/internal/model"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
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

// UpsertBySlug creates or updates a post identified by slug.
// New posts land as unpublished drafts (D-11). On update, existing
// tags are preserved -- the admin sets tags from the web UI.
func (s *Service) UpsertBySlug(ctx context.Context, title, slug, body string) error {
	rendered := s.renderer.Render(body)
	existing, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, postrepo.ErrNotFound) {
			_, createErr := s.repo.Create(ctx, model.Post{
				Title:        title,
				Slug:         slug,
				Body:         body,
				RenderedHTML: string(rendered),
				Published:    false,
			})
			return createErr
		}
		return err
	}
	return s.repo.Update(ctx, model.Post{
		ID:           existing.ID,
		Title:        title,
		Slug:         slug,
		Body:         body,
		RenderedHTML: string(rendered),
		Tags:         existing.Tags,
	})
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
