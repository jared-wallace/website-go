package post

import (
	"context"
	"html/template"

	"github.com/jared-wallace/website-go/internal/model"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
)

// PostDetail is a fully-hydrated post suitable for the single-post view.
// RenderedHTML is typed as template.HTML to prevent double-escaping in
// Go templates (Pitfall 3 from RESEARCH.md).
type PostDetail struct {
	Post         model.Post
	RenderedHTML template.HTML // cast from model.Post.RenderedHTML; safe — stored pre-sanitized
	ToC          []ToCEntry    // nil when fewer than 3 headings
	Tags         []string
	Excerpt      string // plain-text excerpt for OG description
	ReadingTime  int
}

// GetBySlug fetches a single published post by slug and enriches it with
// derived fields. Returns ErrNotFound (from the repository package) when no
// matching post exists — callers should check with errors.Is.
func (s *Service) GetBySlug(ctx context.Context, slug string) (*PostDetail, error) {
	p, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	enrichedHTML := InjectHeadingIDs(p.RenderedHTML)
	return &PostDetail{
		Post:         *p,
		RenderedHTML: template.HTML(enrichedHTML), //nolint:gosec // pre-sanitized by bluemonday at write time
		ToC:          ExtractToC(enrichedHTML),
		Tags:         ParseTags(p.Tags),
		Excerpt:      Excerpt(p.Body, 200),
		ReadingTime:  ReadingTime(p.Body),
	}, nil
}

// ErrNotFound re-exports the repository sentinel for caller convenience
// — handlers import service/post only and should not need to import both packages.
var ErrNotFound = postrepo.ErrNotFound
