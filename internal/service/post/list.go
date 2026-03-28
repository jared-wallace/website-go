package post

import (
	"context"
	"math"
	"time"
)

// PostsPerPage is the default page size for the public post listing. (D-07)
const PostsPerPage = 10

// PostSummary is a lightweight projection of a Post used in listing views.
// It includes computed fields (Excerpt, Tags, ReadingTime) so templates need
// no business logic.
type PostSummary struct {
	Title       string
	Slug        string
	Excerpt     string
	Tags        []string
	PublishedAt time.Time
	ReadingTime int
}

// ListResult carries a page of PostSummary values together with pagination metadata.
type ListResult struct {
	Posts       []PostSummary
	CurrentPage int
	TotalPages  int
	HasPrev     bool
	HasNext     bool
}

// ListPublished fetches a page of published posts and returns them with
// pagination metadata. Page numbers are 1-based and are clamped:
//   - page < 1 → clamped to 1
//   - page > totalPages → clamped to totalPages
func (s *Service) ListPublished(ctx context.Context, page int) (ListResult, error) {
	if page < 1 {
		page = 1
	}

	total, err := s.repo.CountPublished(ctx)
	if err != nil {
		return ListResult{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(PostsPerPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * PostsPerPage
	rawPosts, err := s.repo.ListPublished(ctx, PostsPerPage, offset)
	if err != nil {
		return ListResult{}, err
	}

	summaries := make([]PostSummary, 0, len(rawPosts))
	for _, p := range rawPosts {
		summaries = append(summaries, PostSummary{
			Title:       p.Title,
			Slug:        p.Slug,
			Excerpt:     Excerpt(p.Body, 200),
			Tags:        ParseTags(p.Tags),
			PublishedAt: p.CreatedAt,
			ReadingTime: ReadingTime(p.Body),
		})
	}

	return ListResult{
		Posts:       summaries,
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}, nil
}
