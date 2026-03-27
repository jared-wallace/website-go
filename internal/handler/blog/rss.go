package blog

import (
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/jared-wallace/website-go/internal/model"
	postservice "github.com/jared-wallace/website-go/internal/service/post"
)

// baseURL is the canonical origin for all absolute URLs in RSS and sitemap output.
const baseURL = "https://jared-wallace.com"

// rssManagingEditor is the RFC 2822 formatted editor field for the RSS channel.
const rssManagingEditor = "jaredwallace@jared-wallace.com (Jared Wallace)"

// --- CDATA type ---

// CDATA wraps a raw HTML string so encoding/xml emits a proper CDATA section
// instead of entity-escaping angle brackets. This avoids Pitfall 1 from the
// research notes: RSS readers seeing literal &lt;p&gt; in descriptions.
type CDATA struct {
	Value string
}

// MarshalXML emits <![CDATA[...]]> so RSS readers receive unescaped HTML.
func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		S string `xml:",innerxml"`
	}{S: "<![CDATA[" + c.Value + "]]>"}, start)
}

// --- RSS 2.0 struct types (per RESEARCH.md Pattern 2) ---

// RSSFeed is the top-level RSS 2.0 document.
type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

// RSSChannel holds the feed metadata and all items.
type RSSChannel struct {
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	Description    string    `xml:"description"`
	Language       string    `xml:"language"`
	ManagingEditor string    `xml:"managingEditor"`
	LastBuildDate  string    `xml:"lastBuildDate"`
	Items          []RSSItem `xml:"item"`
}

// RSSItem represents a single post entry in the feed.
type RSSItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description CDATA    `xml:"description"`
	Author      string   `xml:"author"`
	PubDate     string   `xml:"pubDate"`
	GUID        RSSGuid  `xml:"guid"`
	Categories  []string `xml:"category"`
}

// RSSGuid is the unique identifier for a feed item.
type RSSGuid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// buildRSSFeed constructs an RSSFeed from a slice of published posts.
func buildRSSFeed(posts []model.Post) RSSFeed {
	items := make([]RSSItem, 0, len(posts))
	for _, p := range posts {
		postURL := baseURL + "/posts/" + p.Slug
		items = append(items, RSSItem{
			Title:       p.Title,
			Link:        postURL,
			Description: CDATA{Value: p.RenderedHTML},
			Author:      rssManagingEditor,
			PubDate:     p.CreatedAt.Format(time.RFC1123Z),
			GUID:        RSSGuid{IsPermaLink: true, Value: postURL},
			Categories:  postservice.ParseTags(p.Tags),
		})
	}
	return RSSFeed{
		Version: "2.0",
		Channel: RSSChannel{
			Title:          "The Log",
			Link:           baseURL,
			Description:    "dispatches from the deep end",
			Language:       "en-us",
			ManagingEditor: rssManagingEditor,
			LastBuildDate:  time.Now().Format(time.RFC1123Z),
			Items:          items,
		},
	}
}

// ServeRSS handles GET /rss, returning an RSS 2.0 feed of the 25 most recent
// published posts with full HTML content in CDATA-wrapped descriptions.
func (h *BlogHandler) ServeRSS(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.ListForFeed(r.Context(), 25)
	if err != nil {
		slog.Error("ListForFeed failed", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	feed := buildRSSFeed(posts)
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	if _, err := io.WriteString(w, xml.Header); err != nil {
		slog.Error("rss write header failed", "error", err)
		return
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(feed); err != nil {
		slog.Error("rss encode failed", "error", err)
	}
}
