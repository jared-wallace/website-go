# Phase 4: Distribution - Research

**Researched:** 2026-03-27
**Domain:** RSS 2.0, Open Graph / Twitter Card meta tags, XML sitemap, thumbs-up reactions
**Confidence:** HIGH

---

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions

**RSS Feed**
- D-01: Most recent 25 published posts in the feed. Full post content in `<description>`.
- D-02: Author "Jared Wallace" in `<managingEditor>` and per-item `<author>`.
- D-03: Tags mapped to `<category>` elements per RSS item.
- D-04: Auto-discovery `<link rel="alternate" type="application/rss+xml">` added to base.html so every page advertises the feed.

**Open Graph & Social Previews**
- D-05: OG description sourced from the existing post excerpt. No manual override.
- D-06: Static site-wide fallback OG image (branded nautical graphic) served from `/static/`.
- D-07: Twitter card type `summary` (small square + title + description).
- D-08: Homepage gets site-level OG tags: `og:title="The Log"`, `og:description="dispatches from the deep end"`, `og:image=fallback`.
- D-09: Individual post pages get per-post OG tags: `og:title={post title}`, `og:description={post excerpt}`, `og:image=fallback`.

**Thumbs-Up Reactions**
- D-10: Button placed below post content, before footer. Shows current count.
- D-11: One thumbs-up per reader per post (binary, not multi-clap).
- D-12: Rate limiting via IP-based server-side check (one per IP per post per 24h, IP hash + post ID in reactions table) plus localStorage flag client-side.
- D-13: Tap feedback: count increments with subtle CSS animation. No JS animation library.

**Sitemap & Crawlers**
- D-14: Sitemap includes all published post URLs plus homepage URL.
- D-15: `/robots.txt` handler returns `Sitemap: https://jared-wallace.com/sitemap.xml` plus standard `Allow` directive.

### Claude's Discretion
- RSS feed channel-level metadata (use "The Log" / "dispatches from the deep end")
- OG fallback image design (nautical-themed, appropriate dimensions for `summary` card)
- Thumbs-up button icon and CSS animation specifics
- Sitemap `<changefreq>` and `<priority>` values
- Reactions table schema (IP hashing approach, index design)
- Whether to use a `like_count` column on posts table vs. COUNT query on reactions table
- robots.txt additional directives (if any)

### Deferred Ideas (OUT OF SCOPE)
None — discussion stayed within phase scope.
</user_constraints>

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| BLOG-06 | Reader sees proper OG meta tags when links are shared | OG/Twitter Card meta patterns in `<head>` block, per-page template override via new `{{block "meta" .}}` |
| BLOG-07 | Search engines can discover all posts via /sitemap.xml | Sitemap 0.9 XML schema, `encoding/xml` struct pattern, `/robots.txt` handler |
| BLOG-09 | Reader can subscribe via RSS feed at /rss with full post content | RSS 2.0 spec, `encoding/xml` struct pattern, existing `ListPublished` + new feed limit |
| BLOG-10 | Reader can give thumbs-up reaction on posts (rate-limited, no auth required) | Reactions table schema, IP hashing, `POST /posts/{slug}/react` handler, localStorage guard, CSS animation |
</phase_requirements>

---

## Summary

Phase 4 adds four distribution features to the existing blog: RSS 2.0 feed, Open Graph social previews, XML sitemap with robots.txt, and a rate-limited thumbs-up reaction system. The codebase is well-structured for this work — all four features follow established patterns (new handler methods on `BlogHandler`, new routes in `blogMux`, `encoding/xml` for the XML endpoints, one new migration for reactions).

The most substantive new work is the reactions system, which introduces the only new persistent model. It requires a new repository interface, a new migration (00004), and minimal JavaScript added to `main.js`. All other features are essentially new HTTP handlers that write XML or inject data into existing templates.

The primary template change needed is adding a `{{block "meta" .}}` to `base.html`'s `<head>` section — it currently has no such extension point. Every page that needs custom OG tags (list and post pages) will override this block. The RSS auto-discovery link (D-04) goes in the base as a static tag, not in the block.

**Primary recommendation:** Implement in four discrete units: (1) RSS handler, (2) OG meta tags + auto-discovery, (3) sitemap + robots.txt, (4) reactions system. Each unit is independently testable and has no cross-dependencies.

---

## Project Constraints (from CLAUDE.md)

- Go with minimal dependencies — use `encoding/xml` stdlib for RSS and sitemap, not external libraries.
- All persistent data lives on EBS at `/var/www/html` — no new external storage (reactions go in Postgres).
- `frontend-design` skill referenced for template/UI work — no skills directory found in this repo, follow existing nautical CSS patterns in `main.css`.
- Run code changes through `/simplify` before finalizing.
- Entry point for file changes must be a GSD workflow.

---

## Standard Stack

### Core (all stdlib — no new dependencies)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `encoding/xml` | Go 1.26 stdlib | RSS 2.0 feed + sitemap XML generation | Locked in STACK.md — RSS is straightforward XML struct marshaling |
| `net/http` | Go 1.26 stdlib | `/rss`, `/sitemap.xml`, `/robots.txt`, `/posts/{slug}/react` routes | Existing routing pattern |
| `crypto/sha256` | Go 1.26 stdlib | IP address hashing before storing in reactions table | Privacy-preserving, no external dep |
| `html/template` | Go 1.26 stdlib | OG meta block in base.html + post/list template overrides | Already used for all templates |

### No New Dependencies Required
The go.mod already contains everything needed. `encoding/xml` is stdlib. The reactions endpoint uses existing pgx/v5 pool. CSS animations are plain CSS. JavaScript uses `fetch()` and `localStorage` (browser APIs, no library).

---

## Architecture Patterns

### Recommended Project Structure (new files only)
```
internal/
├── handler/blog/
│   ├── rss.go              # ServeRSS handler method
│   ├── sitemap.go          # ServeSitemap + ServeRobots handler methods
│   └── react.go            # React handler method (POST /posts/{slug}/react)
├── model/
│   └── reaction.go         # Reaction struct (new)
├── repository/post/
│   └── reactions.go        # CountReactions, AddReaction, HasReacted queries
└── service/post/
    └── rss.go              # ListForFeed service method (top 25 published)

db/migrations/
└── 00004_create_reactions.sql

web/
├── static/
│   └── og-fallback.png     # Nautical branded OG fallback image (1200×630 for summary_large or 400×400 for summary)
└── templates/
    ├── base.html           # +RSS auto-discovery link, +{{block "meta" .}}{{end}}
    ├── list.html           # +{{define "meta"}} with site-level OG tags
    └── post.html           # +{{define "meta"}} with per-post OG tags, +thumbs-up button
```

### Pattern 1: XML Handler (RSS and Sitemap)
**What:** Handler method writes XML directly to `http.ResponseWriter` using `encoding/xml`. No template involved.
**When to use:** Any endpoint that returns structured XML rather than HTML.
**Example:**
```go
// Source: encoding/xml stdlib docs
func (h *BlogHandler) ServeRSS(w http.ResponseWriter, r *http.Request) {
    posts, err := h.svc.ListForFeed(r.Context(), 25)
    if err != nil {
        slog.Error("ListForFeed failed", "error", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }
    feed := buildRSSFeed(posts)
    w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    enc := xml.NewEncoder(w)
    enc.Indent("", "  ")
    if err := enc.Encode(feed); err != nil {
        slog.Error("rss encode failed", "error", err)
    }
}
```

### Pattern 2: RSS 2.0 XML Struct
**What:** Go struct with `xml` tags matching the RSS 2.0 spec. Marshaled directly with `encoding/xml`.
**When to use:** RSS feed generation.
**Example:**
```go
// Source: RSS 2.0 spec (https://www.rssboard.org/rss-specification)
type RSSFeed struct {
    XMLName xml.Name   `xml:"rss"`
    Version string     `xml:"version,attr"`
    Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
    Title          string    `xml:"title"`
    Link           string    `xml:"link"`
    Description    string    `xml:"description"`
    Language       string    `xml:"language"`
    ManagingEditor string    `xml:"managingEditor"`
    LastBuildDate  string    `xml:"lastBuildDate"`
    Items          []RSSItem `xml:"item"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"` // full HTML content per D-01
    Author      string `xml:"author"`
    PubDate     string `xml:"pubDate"`
    GUID        RSSGuid `xml:"guid"`
    Categories  []string `xml:"category"`
}

type RSSGuid struct {
    IsPermaLink bool   `xml:"isPermaLink,attr"`
    Value       string `xml:",chardata"`
}
```

**RFC 822 date format for RSS:** `time.RFC1123Z` ("Mon, 02 Jan 2006 15:04:05 -0700")

### Pattern 3: Sitemap 0.9 XML Struct
**What:** Go struct matching the sitemaps.org schema. Marshaled with `encoding/xml`.
**When to use:** `/sitemap.xml` endpoint.
**Example:**
```go
// Source: https://www.sitemaps.org/protocol.html
type URLSet struct {
    XMLName xml.Name  `xml:"urlset"`
    XMLNS   string    `xml:"xmlns,attr"`
    URLs    []SitemapURL `xml:"url"`
}
// XMLNS value: "http://www.sitemaps.org/schemas/sitemap/0.9"

type SitemapURL struct {
    Loc        string `xml:"loc"`
    LastMod    string `xml:"lastmod,omitempty"` // YYYY-MM-DD
    ChangeFreq string `xml:"changefreq,omitempty"`
    Priority   string `xml:"priority,omitempty"`
}
```

### Pattern 4: OG Meta Block in base.html
**What:** Add `{{block "meta" .}}{{end}}` to `<head>` in `base.html`, with a default set of site-wide OG tags. Individual pages override with `{{define "meta"}}`.
**When to use:** Extensible per-page OG metadata.

**base.html addition (inside `<head>`):**
```html
<!-- RSS auto-discovery (D-04) — static, always present -->
<link rel="alternate" type="application/rss+xml" title="The Log" href="/rss">
<!-- OG meta block — overridden per page -->
{{block "meta" .}}
<meta property="og:title" content="The Log">
<meta property="og:description" content="dispatches from the deep end">
<meta property="og:type" content="website">
<meta property="og:image" content="https://jared-wallace.com/static/og-fallback.png">
<meta property="og:url" content="https://jared-wallace.com">
<meta name="twitter:card" content="summary">
<meta name="twitter:title" content="The Log">
<meta name="twitter:description" content="dispatches from the deep end">
<meta name="twitter:image" content="https://jared-wallace.com/static/og-fallback.png">
{{end}}
```

**post.html addition (`{{define "meta"}}`):**
```html
{{define "meta"}}
<meta property="og:title" content="{{.Post.Title}}">
<meta property="og:description" content="{{.Excerpt}}">
<meta property="og:type" content="article">
<meta property="og:image" content="https://jared-wallace.com/static/og-fallback.png">
<meta property="og:url" content="https://jared-wallace.com/posts/{{.Post.Slug}}">
<meta name="twitter:card" content="summary">
<meta name="twitter:title" content="{{.Post.Title}}">
<meta name="twitter:description" content="{{.Excerpt}}">
<meta name="twitter:image" content="https://jared-wallace.com/static/og-fallback.png">
{{end}}
```

**CRITICAL NOTE on per-page template sets:** The existing `New()` parses base.html + page.html per page into separate `*template.Template` instances using `template.Must(...ParseFS(..., "base.html", page))`. The `{{define "meta"}}` in post.html is parsed into the same template set as base.html, so the override works. No architecture change needed.

### Pattern 5: Reactions Table Schema
**What:** New table for tracking per-reader reactions per post. IP stored as SHA-256 hash for privacy.
**Recommendation:** Use a `COUNT(*)` query on reactions table rather than a `like_count` denormalized column. At blog scale (~hundreds of posts, ~thousands of reactions), the COUNT query with a proper index is fast. Avoids write-on-read (updating a column on every reaction) and keeps the schema clean.

**Migration SQL:**
```sql
-- +goose Up
CREATE TABLE reactions (
    id         BIGSERIAL   PRIMARY KEY,
    post_id    BIGINT      NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    ip_hash    TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- One reaction per IP per post enforced at DB level
CREATE UNIQUE INDEX reactions_post_ip_uidx ON reactions (post_id, ip_hash);
-- Fast count lookup by post
CREATE INDEX reactions_post_id_idx ON reactions (post_id);

-- +goose Down
DROP TABLE IF EXISTS reactions;
```

**IP Hashing:**
```go
// Source: crypto/sha256 stdlib
import "crypto/sha256"
import "fmt"

func hashIP(ip string) string {
    sum := sha256.Sum256([]byte(ip))
    return fmt.Sprintf("%x", sum)
}
```

Extract client IP with care (see Pitfall 3):
```go
func clientIP(r *http.Request) string {
    // Behind Nginx ALB — trust X-Real-IP set by Nginx
    if ip := r.Header.Get("X-Real-IP"); ip != "" {
        return ip
    }
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

### Pattern 6: React Endpoint (POST /posts/{slug}/react)
**What:** Inserts a reaction row. Returns JSON `{"count": N, "already_reacted": bool}`. Uses `ON CONFLICT DO NOTHING` to handle the unique constraint gracefully.
**When to use:** Thumbs-up button fetch() call.

```go
// Response JSON only — no template rendering
type reactResponse struct {
    Count          int  `json:"count"`
    AlreadyReacted bool `json:"already_reacted"`
}
```

**Repository methods needed:**
```go
// AddReaction inserts a reaction. Returns (alreadyExists bool, err error).
// Uses INSERT ... ON CONFLICT (post_id, ip_hash) DO NOTHING.
AddReaction(ctx context.Context, postID int64, ipHash string) (bool, error)

// CountReactions returns the reaction count for a post.
CountReactions(ctx context.Context, postID int64) (int, error)
```

### Pattern 7: Thumbs-Up Button HTML + JavaScript
**What:** Button in post.html, below post body. JavaScript in main.js handles the fetch() and localStorage guard.

**post.html addition (after `.post-body`, before `{{end}}`):**
```html
<div class="reaction-bar">
  <button class="reaction-btn" id="reaction-btn" data-slug="{{.Post.Slug}}" aria-label="Give a thumbs-up">
    <span class="reaction-icon">&#128077;</span>
    <span class="reaction-count" id="reaction-count">{{.ReactionCount}}</span>
  </button>
</div>
```

**main.js addition:**
```javascript
// Thumbs-up reaction
var reactionBtn = document.getElementById('reaction-btn');
if (reactionBtn) {
  var slug = reactionBtn.getAttribute('data-slug');
  var storageKey = 'reacted:' + slug;

  // Disable button if already reacted (localStorage guard)
  if (localStorage.getItem(storageKey)) {
    reactionBtn.disabled = true;
    reactionBtn.classList.add('reacted');
  }

  reactionBtn.addEventListener('click', function() {
    if (reactionBtn.disabled) return;
    reactionBtn.disabled = true;

    fetch('/posts/' + slug + '/react', { method: 'POST' })
      .then(function(res) { return res.json(); })
      .then(function(data) {
        document.getElementById('reaction-count').textContent = data.count;
        reactionBtn.classList.add('reacted');
        localStorage.setItem(storageKey, '1');
      })
      .catch(function() {
        // Silently fail — don't punish the reader for a network hiccup
        reactionBtn.disabled = false;
      });
  });
}
```

### Pattern 8: robots.txt Handler
**What:** Plain text response — no template, no XML. Inline string write.

```go
func (h *BlogHandler) ServeRobots(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    fmt.Fprintf(w, "User-agent: *\nAllow: /\nSitemap: https://jared-wallace.com/sitemap.xml\n")
}
```

### Anti-Patterns to Avoid
- **Storing raw IP addresses in reactions table:** Always SHA-256 hash before storing. GDPR/privacy hygiene.
- **Using `X-Forwarded-For` to extract client IP:** This header can be spoofed. Use `X-Real-IP` (set by Nginx) or fall back to `RemoteAddr`. Do not trust `X-Forwarded-For` for rate limiting.
- **Returning 200 on duplicate reaction:** Return 200 with `already_reacted: true` — don't 409 the user. The localStorage guard handles the UX; the server just confirms state.
- **`{{define "meta"}}` in a separate template set:** Because the existing pattern parses base.html + page.html into the same template set per page, defining "meta" in the page template works naturally. Don't create a third template file.
- **Storing full HTML in RSS `<description>`:** Use `<![CDATA[...]]>` wrapping when using `encoding/xml` to include raw HTML without escaping. In Go, use `xml.CharData` or a wrapper type. (See Pitfall 1.)

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| RSS XML generation | Custom string builder | `encoding/xml` struct + Marshal | Handles escaping, attribute quoting, nesting; already in stdlib |
| Sitemap XML | Template-rendered XML | `encoding/xml` struct + Marshal | Same rationale; sitemaps.org schema is small and fits a struct cleanly |
| IP rate limiting (24h window) | Time-based map in memory | Reactions table with `UNIQUE INDEX (post_id, ip_hash)` + timestamp filter | Survives container restarts; consistent across multiple concurrent requests |
| OG image resizing | Dynamic image processing | Static PNG served from `/static/` | Per D-06; no need for image processing at this phase |
| CSS bounce animation | JS animation library | `@keyframes` in main.css | Simple one-shot animation; no library needed |

**Key insight:** All four features in this phase have known, simple implementations using stdlib. The complexity ceiling is low — resist the temptation to reach for external packages.

---

## Common Pitfalls

### Pitfall 1: HTML in RSS `<description>` gets double-escaped
**What goes wrong:** `encoding/xml` escapes `<` and `>` in string fields to `&lt;` and `&gt;`. RSS readers display literal HTML tags instead of rendered content.
**Why it happens:** Go's `encoding/xml` treats string fields as text content and escapes all XML special characters.
**How to avoid:** Wrap the HTML content in a custom type that implements `xml.Marshaler` and emits a CDATA section:
```go
type CDATA struct {
    Value string
}

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return e.EncodeElement(struct {
        S string `xml:",innerxml"`
    }{S: "<![CDATA[" + c.Value + "]]>"}, start)
}
```
Use `CDATA{Value: post.RenderedHTML}` in the `RSSItem.Description` field.
**Warning signs:** RSS feed validator reports entity-encoded HTML; feed readers show `&lt;p&gt;` in description.

### Pitfall 2: `{{block "meta"}}` override doesn't render when template set is reused
**What goes wrong:** `{{define "meta"}}` in post.html is silently ignored; the default from base.html renders instead.
**Why it happens:** In `html/template`, a `{{define}}` block in a file must be parsed into the same `*template.Template` instance as the `{{block}}` that references it. The existing pattern already does this correctly (base.html + page.html parsed together), so this pitfall is avoided as long as the pattern is followed.
**How to avoid:** Verify `handler.go`'s `New()` parses `base.html` + `post.html` together into the same set. Do not split them.
**Warning signs:** Post pages show site-level OG title even after adding `{{define "meta"}}` to post.html.

### Pitfall 3: IP extraction behind Nginx reverse proxy gives `127.0.0.1` for all requests
**What goes wrong:** `r.RemoteAddr` is always `127.0.0.1` (the Nginx container IP). All readers share one "IP" and rate limiting is broken.
**Why it happens:** The existing Nginx config sets `proxy_pass` to the Go server. Without extracting the real client IP from a trusted header, `r.RemoteAddr` is the proxy's address.
**How to avoid:** Use `X-Real-IP` header (set by Nginx's `proxy_set_header X-Real-IP $remote_addr`). Verify this header is already set in the Nginx config. Fall back to `r.RemoteAddr` only in local dev.
**Warning signs:** In production, all reactions appear to come from the same IP; rate limiting fires for the first reader and blocks everyone else.

### Pitfall 4: Twitter `summary` card image must be at least 144×144px
**What goes wrong:** Twitter/X silently ignores the OG image and shows no image in the card.
**Why it happens:** Twitter `summary` card type has minimum image size requirements. The image must be at least 144×144px and ideally square (Twitter crops to 1:1 for `summary`).
**How to avoid:** Create the fallback OG image at exactly 400×400px (recommended for `summary` card). This is within Twitter's 5MB limit and above the minimum.
**Warning signs:** Twitter Card Validator shows no image preview; card appears without image.

### Pitfall 5: Sitemap `<loc>` URLs must be absolute with scheme
**What goes wrong:** Relative URLs like `/posts/my-post` in the sitemap are invalid per the sitemaps.org spec. Google Search Console rejects the sitemap.
**Why it happens:** Developers write slug-only paths, forgetting the full base URL.
**How to avoid:** Always prefix with `https://jared-wallace.com` in the sitemap handler. Store the base URL as a constant or inject from config.
**Warning signs:** Google Search Console reports "Invalid URL" errors when processing the sitemap.

### Pitfall 6: Reaction button re-fires if fetch is slow (double-tap)
**What goes wrong:** User taps thumbs-up, fetch is slow, user taps again before response — two requests fire, second one hits the UNIQUE constraint and errors.
**Why it happens:** Button is only disabled after the response arrives in the `.then()` handler, not immediately on click.
**How to avoid:** Disable the button immediately in the click handler (before the fetch), re-enable only on error. The code example in Pattern 7 above does this correctly.
**Warning signs:** Occasional 500 errors from the react endpoint in logs; double-counted reactions.

---

## Code Examples

### RSS Channel Build Function
```go
// Source: encoding/xml stdlib + RSS 2.0 spec
const baseURL = "https://jared-wallace.com"

func buildRSSFeed(posts []PostFeedItem) RSSFeed {
    items := make([]RSSItem, 0, len(posts))
    for _, p := range posts {
        postURL := baseURL + "/posts/" + p.Slug
        items = append(items, RSSItem{
            Title:       p.Title,
            Link:        postURL,
            Description: CDATA{Value: string(p.RenderedHTML)},
            Author:      "jaredwallace@example.com (Jared Wallace)",
            PubDate:     p.PublishedAt.Format(time.RFC1123Z),
            GUID:        RSSGuid{IsPermaLink: true, Value: postURL},
            Categories:  p.Tags,
        })
    }
    return RSSFeed{
        Version: "2.0",
        Channel: RSSChannel{
            Title:          "The Log",
            Link:           baseURL,
            Description:    "dispatches from the deep end",
            Language:       "en-us",
            ManagingEditor: "jaredwallace@example.com (Jared Wallace)",
            LastBuildDate:  time.Now().Format(time.RFC1123Z),
            Items:          items,
        },
    }
}
```

### Sitemap Build Function
```go
// Source: https://www.sitemaps.org/protocol.html
func buildSitemap(posts []PostSitemapItem) URLSet {
    urls := make([]SitemapURL, 0, len(posts)+1)
    // Homepage
    urls = append(urls, SitemapURL{
        Loc:        "https://jared-wallace.com/",
        ChangeFreq: "daily",
        Priority:   "1.0",
    })
    for _, p := range posts {
        urls = append(urls, SitemapURL{
            Loc:        "https://jared-wallace.com/posts/" + p.Slug,
            LastMod:    p.PublishedAt.Format("2006-01-02"),
            ChangeFreq: "monthly",
            Priority:   "0.8",
        })
    }
    return URLSet{
        XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
        URLs:  urls,
    }
}
```

### ADD XML Declaration (Required for RSS and Sitemap)
Both RSS and sitemap should emit an XML declaration. Use `io.WriteString` before the encoder:
```go
io.WriteString(w, xml.Header) // writes: <?xml version="1.0" encoding="UTF-8"?>
enc := xml.NewEncoder(w)
```

### Reaction INSERT with Conflict Handling
```go
// Source: pgx v5 docs, PostgreSQL ON CONFLICT syntax
const insertReaction = `
    INSERT INTO reactions (post_id, ip_hash)
    VALUES ($1, $2)
    ON CONFLICT (post_id, ip_hash) DO NOTHING`

func (r *postgresRepository) AddReaction(ctx context.Context, postID int64, ipHash string) (alreadyExists bool, err error) {
    tag, err := r.pool.Exec(ctx, insertReaction, postID, ipHash)
    if err != nil {
        return false, fmt.Errorf("AddReaction: %w", err)
    }
    return tag.RowsAffected() == 0, nil
}
```

### CSS Bounce Animation for Reaction
```css
/* In main.css */
.reaction-btn {
  background: none;
  border: 2px solid var(--accent);
  border-radius: 2rem;
  padding: 0.4rem 1rem;
  cursor: pointer;
  font-size: 1rem;
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  transition: background 0.2s;
}

.reaction-btn.reacted,
.reaction-btn:disabled {
  background: var(--accent);
  cursor: default;
}

.reaction-btn.bounce .reaction-icon {
  animation: reaction-bounce 0.4s ease;
}

@keyframes reaction-bounce {
  0%   { transform: scale(1); }
  40%  { transform: scale(1.4); }
  70%  { transform: scale(0.9); }
  100% { transform: scale(1); }
}
```

Add `reactionBtn.classList.add('bounce')` in the fetch `.then()` handler; remove it after the animation ends (`animationend` event) for re-trigger safety.

---

## Service Layer: New Method Needed

The existing `ListPublished` is paginated (10 per page). RSS needs the top 25, and sitemap needs all published posts. A new service method (or a new repository query) is required.

**Option A (recommended):** New `ListForFeed(ctx, limit int)` service method that calls a new `ListPublishedLimit(ctx, limit)` repository method. This is a simple `LIMIT $1 OFFSET 0` query, reusing the existing query pattern.

**Option B:** Reuse `repo.ListPublished(ctx, 25, 0)` directly from the handler. Simpler, but bypasses the service layer.

**Recommendation:** Option A, for consistency with the existing pattern (handlers never call repo directly). The method is trivial to add and keeps the architecture clean.

For sitemap, a `ListAllPublishedForSitemap(ctx)` method returning only `(slug, created_at)` tuples avoids fetching body/rendered_html for all posts just to build URLs.

---

## Route Registration (main.go additions)

The following routes are added to `blogMux` in `cmd/server/main.go`:

```go
blogMux.HandleFunc("GET /rss", blog.ServeRSS)
blogMux.HandleFunc("GET /sitemap.xml", blog.ServeSitemap)
blogMux.HandleFunc("GET /robots.txt", blog.ServeRobots)
blogMux.HandleFunc("POST /posts/{slug}/react", blog.React)
```

These must be registered **before** the catch-all `GET /{path...}` handler.

---

## Environment Availability

Step 2.6: SKIPPED (no external dependencies — all work is code/config changes using existing stack)

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` + `net/http/httptest` |
| Config file | none (standard `go test ./...`) |
| Quick run command | `go test ./internal/handler/blog/... ./internal/service/post/... ./internal/repository/post/...` |
| Full suite command | `go test ./...` |

### Phase Requirements → Test Map

| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|--------------|
| BLOG-09 | `/rss` returns valid RSS 2.0 XML with only published posts | unit | `go test ./internal/handler/blog/... -run TestServeRSS` | ❌ Wave 0 |
| BLOG-09 | RSS feed contains no draft posts | unit | `go test ./internal/handler/blog/... -run TestRSSDraftExclusion` | ❌ Wave 0 |
| BLOG-09 | RSS feed contains full post content in description | unit | `go test ./internal/handler/blog/... -run TestRSSFullContent` | ❌ Wave 0 |
| BLOG-07 | `/sitemap.xml` lists all published post URLs | unit | `go test ./internal/handler/blog/... -run TestServeSitemap` | ❌ Wave 0 |
| BLOG-07 | `/robots.txt` includes Sitemap directive | unit | `go test ./internal/handler/blog/... -run TestServeRobots` | ❌ Wave 0 |
| BLOG-06 | Post page HTML contains correct `og:title` meta tag | unit | `go test ./internal/handler/blog/... -run TestPostOGMeta` | ❌ Wave 0 |
| BLOG-06 | Homepage HTML contains site-level OG tags | unit | `go test ./internal/handler/blog/... -run TestListOGMeta` | ❌ Wave 0 |
| BLOG-10 | `POST /posts/{slug}/react` increments count | unit | `go test ./internal/handler/blog/... -run TestReact` | ❌ Wave 0 |
| BLOG-10 | Duplicate reaction from same IP returns `already_reacted: true` | unit | `go test ./internal/handler/blog/... -run TestReactDuplicate` | ❌ Wave 0 |
| BLOG-10 | Reaction count in post page template | unit | `go test ./internal/handler/blog/... -run TestPostReactionCount` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./internal/handler/blog/... -run TestServe`
- **Per wave merge:** `go test ./...`
- **Phase gate:** Full suite green before `/gsd:verify-work`

### Wave 0 Gaps
- [ ] `internal/handler/blog/rss_test.go` — covers BLOG-09 RSS tests
- [ ] `internal/handler/blog/sitemap_test.go` — covers BLOG-07
- [ ] `internal/handler/blog/react_test.go` — covers BLOG-10
- [ ] `internal/repository/post/reactions_test.go` — covers AddReaction and CountReactions (mock-level, no DB required for unit tests)
- [ ] The existing `mockRepository` in `handler_test.go` needs two new interface methods: `AddReaction` and `CountReactions` — add to existing mock

**Note:** OG meta tests (BLOG-06) can be covered in the existing `handler_test.go` by asserting on the rendered HTML body from `ListPosts` and `ShowPost`. No new test file needed for these.

---

## Open Questions

1. **Does Nginx already set `X-Real-IP`?**
   - What we know: The app is behind Nginx. The reactions rate limiting depends on accurate client IP extraction.
   - What's unclear: Whether the existing Nginx config sets `proxy_set_header X-Real-IP $remote_addr`.
   - Recommendation: The planner should include a task to verify the Nginx config and add the header if missing, OR document that the reactions handler uses `X-Real-IP` with a fallback — and note that rate limiting will be degraded (all share one IP) if the header is absent. For a personal blog this is low risk.

2. **`managingEditor` email address for RSS**
   - What we know: D-02 specifies "Jared Wallace" as author. RSS 2.0 spec requires `managingEditor` to be an email address (format: `email (Name)`).
   - What's unclear: Whether the author's real email should be in the feed (public) or a placeholder.
   - Recommendation: Use `jaredwallace@jared-wallace.com (Jared Wallace)` or omit `managingEditor` entirely (optional field). Claude's discretion per CONTEXT.md.

3. **Static OG image creation**
   - What we know: D-06 requires a branded nautical graphic. The `frontend-design` skill is referenced in CLAUDE.md but no skills directory exists in this repo.
   - What's unclear: Whether the planner should create the PNG as part of this phase or treat it as a placeholder.
   - Recommendation: Create a simple placeholder PNG (400×400, solid color with "The Log" text) in Wave 0. A real branded image can replace it without any code changes.

---

## Sources

### Primary (HIGH confidence)
- Go stdlib `encoding/xml` — struct + Marshal pattern, CDATA handling
- RSS 2.0 specification: https://www.rssboard.org/rss-specification — field names, date format, GUID semantics
- Sitemaps.org protocol: https://www.sitemaps.org/protocol.html — `<urlset>`, `<url>`, `<loc>` requirements
- Existing codebase (direct inspection): handler patterns, template set parsing, repository interface, model structs, go.mod dependency versions

### Secondary (MEDIUM confidence)
- Twitter Card validator docs — `summary` card image minimum size (144×144, recommended 400×400)
- Open Graph protocol: https://ogp.me — `og:type`, `og:title`, `og:description`, `og:image` properties
- PostgreSQL `ON CONFLICT DO NOTHING` — pgx v5 `RowsAffected()` pattern

### Tertiary (LOW confidence)
- None — all claims are verifiable from official specs or the codebase.

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — no new dependencies; all stdlib
- Architecture: HIGH — follows existing codebase patterns directly observed
- RSS/Sitemap XML patterns: HIGH — official spec + stdlib docs
- OG/Twitter Card spec: MEDIUM — spec verified, image size requirements from Twitter docs
- Pitfalls: HIGH for IP extraction / template override / CDATA; MEDIUM for Twitter image minimum size

**Research date:** 2026-03-27
**Valid until:** 2026-06-01 (stable specs; Go stdlib does not change)
