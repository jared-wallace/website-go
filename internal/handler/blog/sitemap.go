package blog

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jared-wallace/website-go/internal/model"
)

// --- Sitemap 0.9 struct types (per RESEARCH.md Pattern 3) ---

// URLSet is the top-level sitemap document element.
type URLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL represents a single URL entry in the sitemap.
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// buildSitemap constructs a URLSet from all published posts, prepending the homepage.
func buildSitemap(posts []model.Post) URLSet {
	urls := make([]SitemapURL, 0, len(posts)+1)
	// Homepage is always first and highest priority (D-14)
	urls = append(urls, SitemapURL{
		Loc:        baseURL + "/",
		ChangeFreq: "daily",
		Priority:   "1.0",
	})
	for _, p := range posts {
		urls = append(urls, SitemapURL{
			Loc:        baseURL + "/posts/" + p.Slug,
			LastMod:    p.CreatedAt.Format("2006-01-02"),
			ChangeFreq: "monthly",
			Priority:   "0.8",
		})
	}
	return URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}
}

// ServeSitemap handles GET /sitemap.xml, returning a Sitemap 0.9 XML document
// listing the homepage and all published post URLs.
func (h *BlogHandler) ServeSitemap(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.ListSlugsForSitemap(r.Context())
	if err != nil {
		slog.Error("ListSlugsForSitemap failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	sitemap := buildSitemap(posts)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	if _, err := io.WriteString(w, xml.Header); err != nil {
		slog.Error("sitemap write header failed", "error", err)
		return
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(sitemap); err != nil {
		slog.Error("sitemap encode failed", "error", err)
	}
}

// ServeRobots handles GET /robots.txt (per RESEARCH.md Pattern 8).
// Returns plain text with the standard Allow directive and Sitemap location (D-15).
func (h *BlogHandler) ServeRobots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "User-agent: *\nAllow: /\nSitemap: %s/sitemap.xml\n", baseURL)
}
