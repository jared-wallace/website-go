// Package post provides the service layer for blog post business logic.
// It sits between HTTP handlers and the repository, computing derived fields
// such as reading time, table of contents, and excerpt.
package post

import (
	"context"
	"html/template"

	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
)

// Renderer is the interface the service uses to convert markdown to HTML.
// The concrete implementation is *markdown.Renderer; the interface allows
// substitution in tests.
type Renderer interface {
	Render(src string) template.HTML
}

// Service coordinates all blog-post operations. It depends on a Repository
// interface so unit tests can substitute a mock without a real database.
type Service struct {
	repo     postrepo.Repository
	renderer Renderer
}

// New creates a Service with the supplied repository and markdown renderer.
func New(repo postrepo.Repository, renderer Renderer) *Service {
	return &Service{repo: repo, renderer: renderer}
}

// newServiceWithRenderer is an internal alias used in tests within this package.
func newServiceWithRenderer(repo postrepo.Repository, renderer Renderer) *Service {
	return New(repo, renderer)
}

// AddReaction records a reader's thumbs-up reaction on a post.
func (s *Service) AddReaction(ctx context.Context, postID int64, ipHash string) (bool, error) {
	return s.repo.AddReaction(ctx, postID, ipHash)
}

// CountReactions returns the total reaction count for a post.
func (s *Service) CountReactions(ctx context.Context, postID int64) (int, error) {
	return s.repo.CountReactions(ctx, postID)
}
