// Package post provides the service layer for blog post business logic.
// It sits between HTTP handlers and the repository, computing derived fields
// such as reading time, table of contents, and excerpt.
package post

import (
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
)

// Service coordinates all blog-post operations. It depends on a Repository
// interface so unit tests can substitute a mock without a real database.
type Service struct {
	repo postrepo.Repository
}

// New creates a Service with the supplied repository.
func New(repo postrepo.Repository) *Service {
	return &Service{repo: repo}
}
