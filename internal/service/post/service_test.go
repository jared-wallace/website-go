package post_test

import (
	"context"
	"errors"
	"html/template"
	"strings"
	"testing"

	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// noopRenderer satisfies postservice.Renderer for tests that never write posts.
type noopRenderer struct{}

func (noopRenderer) Render(src string) template.HTML { return template.HTML(src) }

// ---- ReadingTime ---------------------------------------------------------

func TestReadingTime(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"empty returns minimum 1", "", 1},
		{"exactly 200 words returns 1", strings.Repeat("word ", 200), 1},
		{"201 words returns 2", strings.Repeat("word ", 201), 2},
		{"1000 words returns 5", strings.Repeat("word ", 1000), 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := postservice.ReadingTime(tc.body)
			if got != tc.want {
				t.Errorf("ReadingTime(%d words) = %d, want %d",
					len(strings.Fields(tc.body)), got, tc.want)
			}
		})
	}
}

// ---- ExtractToC ----------------------------------------------------------

func TestExtractToC_FewHeadings(t *testing.T) {
	// 0 headings → nil
	if toc := postservice.ExtractToC("<p>No headings here.</p>"); toc != nil {
		t.Errorf("expected nil for 0 headings, got %v", toc)
	}
	// 2 headings → nil (threshold is 3)
	html2 := `<h2 id="a">Alpha</h2><h2 id="b">Beta</h2>`
	if toc := postservice.ExtractToC(html2); toc != nil {
		t.Errorf("expected nil for 2 headings, got %v", toc)
	}
}

func TestExtractToC_ThreeOrMore(t *testing.T) {
	html3 := `<h2 id="one">One</h2><h3 id="two">Two</h3><h2 id="three">Three</h2>`
	toc := postservice.ExtractToC(html3)
	if toc == nil {
		t.Fatal("expected non-nil ToC for 3 headings")
	}
	if len(toc) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(toc))
	}
}

func TestExtractToC_IDs(t *testing.T) {
	rawHTML := `<h2 id="intro">Introduction</h2><h3 id="sub">Sub-section</h3><h2 id="conc">Conclusion</h2>`
	toc := postservice.ExtractToC(rawHTML)
	if len(toc) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(toc))
	}
	if toc[0].ID != "intro" || toc[0].Text != "Introduction" || toc[0].Level != 2 {
		t.Errorf("entry 0 = %+v, want {ID:intro Text:Introduction Level:2}", toc[0])
	}
	if toc[1].ID != "sub" || toc[1].Text != "Sub-section" || toc[1].Level != 3 {
		t.Errorf("entry 1 = %+v, want {ID:sub Text:Sub-section Level:3}", toc[1])
	}
}

func TestExtractToC_IgnoresH1H4(t *testing.T) {
	html := `<h1 id="h1">Heading1</h1><h2 id="h2">Heading2</h2><h3 id="h3">Heading3</h3><h4 id="h4">Heading4</h4><h2 id="h2b">Second H2</h2>`
	toc := postservice.ExtractToC(html)
	for _, e := range toc {
		if e.Level == 1 || e.Level == 4 {
			t.Errorf("ToC should not include h%d, got %+v", e.Level, e)
		}
	}
	// Should have h2, h3, h2 = 3 entries
	if len(toc) != 3 {
		t.Errorf("expected 3 entries (h2, h3, h2), got %d", len(toc))
	}
}

// ---- Excerpt -------------------------------------------------------------

func TestExcerpt_Short(t *testing.T) {
	body := "This is a short post."
	got := postservice.Excerpt(body, 200)
	if got != body {
		t.Errorf("Excerpt short = %q, want %q", got, body)
	}
}

func TestExcerpt_Truncates(t *testing.T) {
	// Build a body longer than 200 chars
	body := strings.Repeat("word ", 60) // 300 chars
	got := postservice.Excerpt(body, 200)
	if len(got) > 203 { // 200 + "..." is the max (roughly)
		t.Errorf("Excerpt too long: %d chars", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("Excerpt should end with '...', got %q", got)
	}
}

func TestExcerpt_StripsMarkdown(t *testing.T) {
	body := "**bold** and *italic* and [link text](http://example.com) and `code`"
	got := postservice.Excerpt(body, 200)
	for _, unwanted := range []string{"**", "*", "[", "]", "(", ")", "`"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("Excerpt contains markdown syntax %q: %q", unwanted, got)
		}
	}
	if !strings.Contains(got, "bold") || !strings.Contains(got, "italic") || !strings.Contains(got, "link text") || !strings.Contains(got, "code") {
		t.Errorf("Excerpt stripped content it should have kept: %q", got)
	}
}

// ---- ParseTags -----------------------------------------------------------

func TestParseTags(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{"comma-separated", "go,rust,python", []string{"go", "rust", "python"}},
		{"empty string", "", []string{}},
		{"whitespace around commas", " go , rust , python ", []string{"go", "rust", "python"}},
		{"single tag", "go", []string{"go"}},
		{"trailing comma", "go,", []string{"go"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := postservice.ParseTags(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("ParseTags(%q) = %v (len %d), want %v (len %d)",
					tc.input, got, len(got), tc.want, len(tc.want))
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("ParseTags(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
				}
			}
		})
	}
}

// ---- ListPublished pagination --------------------------------------------

// mockRepository implements post.Repository for testing pagination math without a DB.
type mockRepository struct {
	total int
	posts []model.Post
}

func (m *mockRepository) ListPublished(_ context.Context, limit, offset int) ([]model.Post, error) {
	end := offset + limit
	if end > len(m.posts) {
		end = len(m.posts)
	}
	if offset >= len(m.posts) {
		return nil, nil
	}
	return m.posts[offset:end], nil
}

func (m *mockRepository) CountPublished(_ context.Context) (int, error) {
	return m.total, nil
}

func (m *mockRepository) FindBySlug(_ context.Context, slug string) (*model.Post, error) {
	for i := range m.posts {
		if m.posts[i].Slug == slug {
			return &m.posts[i], nil
		}
	}
	return nil, nil
}

func (m *mockRepository) FindByID(_ context.Context, _ int64) (*model.Post, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRepository) ListAll(_ context.Context) ([]model.Post, error) { return nil, nil }

func (m *mockRepository) Create(_ context.Context, _ model.Post) (*model.Post, error) {
	return nil, errors.New("not implemented")
}

func (m *mockRepository) Update(_ context.Context, _ model.Post) error {
	return errors.New("not implemented")
}

func (m *mockRepository) SoftDelete(_ context.Context, _ int64) error {
	return errors.New("not implemented")
}

func (m *mockRepository) Restore(_ context.Context, _ int64) error {
	return errors.New("not implemented")
}

func (m *mockRepository) SetPublished(_ context.Context, _ int64, _ bool) error {
	return errors.New("not implemented")
}

func (m *mockRepository) AddReaction(_ context.Context, _ int64, _ string) (bool, error) {
	return false, nil
}

func (m *mockRepository) CountReactions(_ context.Context, _ int64) (int, error) {
	return 0, nil
}

func makePosts(n int) []model.Post {
	posts := make([]model.Post, n)
	for i := range posts {
		posts[i] = model.Post{ID: int64(i + 1), Slug: "post", Published: true}
	}
	return posts
}

func TestListPublished_Pagination(t *testing.T) {
	repo := &mockRepository{total: 25, posts: makePosts(25)}
	svc := postservice.New(repo, noopRenderer{})

	// Page 1: HasNext=true, HasPrev=false
	r1, err := svc.ListPublished(context.Background(), 1)
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if r1.TotalPages != 3 {
		t.Errorf("page 1 TotalPages = %d, want 3", r1.TotalPages)
	}
	if r1.HasNext != true {
		t.Errorf("page 1 HasNext = false, want true")
	}
	if r1.HasPrev != false {
		t.Errorf("page 1 HasPrev = true, want false")
	}
	if r1.CurrentPage != 1 {
		t.Errorf("page 1 CurrentPage = %d, want 1", r1.CurrentPage)
	}

	// Page 3: HasNext=false, HasPrev=true
	r3, err := svc.ListPublished(context.Background(), 3)
	if err != nil {
		t.Fatalf("page 3: %v", err)
	}
	if r3.HasNext != false {
		t.Errorf("page 3 HasNext = true, want false")
	}
	if r3.HasPrev != true {
		t.Errorf("page 3 HasPrev = false, want true")
	}

	// Page 0 normalizes to page 1
	r0, err := svc.ListPublished(context.Background(), 0)
	if err != nil {
		t.Fatalf("page 0: %v", err)
	}
	if r0.CurrentPage != 1 {
		t.Errorf("page 0 normalized to %d, want 1", r0.CurrentPage)
	}
}
