# Feature Landscape

**Domain:** Personal single-author blog with markdown support (Go)
**Project:** website-go / jared-wallace.com
**Researched:** 2026-03-26
**Overall confidence:** HIGH — established domain with clear community consensus

---

## Table Stakes

Features readers and authors expect. Missing or broken = credibility problem or immediate frustration.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Markdown rendering to HTML | Core contract of a markdown blog | Low | Use goldmark (CommonMark-compliant, extensible, actively maintained). Blackfriday is faster but not CommonMark-compliant and no longer the Go community default. |
| Syntax-highlighted code blocks | Expected by technical readers; ugly unhighlighted code signals low effort | Low | `goldmark-highlighting/v2` wraps Chroma — drop-in extension. Wide theme selection. |
| Image support in posts | Posts without embedded images feel skeletal | Low | Standard markdown `![alt](url)` — store images on EBS volume under `/var/www/html/uploads/`. Serve as static files. |
| RSS feed (RSS 2.0 or Atom) | Feed readers and RSS-first users will check for `/feed` or `/rss`; missing = missed audience | Medium | RSS 2.0 spec is stable. Render full content in `<description>` or summary — readers prefer full-content feeds. |
| Pagination or post listing | Without this, an archive of 50+ posts becomes unusable | Low | Simple offset/limit pagination. 10–15 posts per page is standard. |
| Post published date | Readers use date to assess freshness; required for RSS `<pubDate>` | Low | Store in DB; surface prominently on post and listing. |
| Readable URL slugs | `/posts/my-great-post` vs `/posts/42` — slugs are expected everywhere | Low | Derive from title at create time; allow override; store in DB. |
| Mobile-responsive layout | Google mobile-first indexing; majority of readers arrive on mobile | Medium | CSS-only, no JS required for basic responsiveness. |
| 404 page | Broken links are inevitable; a blank or crash page is embarrassing | Low | Custom handler returning 404 status and a themed error page. |
| Open Graph / social meta tags | Links shared on Slack, Twitter/X, Discord show unfurled previews; missing OG tags = plain URL | Low | `og:title`, `og:description`, `og:image` per post. Cover image field on post or a default site image. |
| Session-based admin auth | Without auth, the admin panel is public; bcrypt + sessions is the correct single-user approach | Medium | Single admin user, bcrypt-hashed password stored in DB or env. Secure session cookie. |
| Web-based markdown editor with preview | Writing in the browser needs live preview to be usable; raw textarea with no preview is painful | High | Split-pane editor: raw markdown left, rendered preview right. Preview can be server-rendered via HTMX partial or client-side JS. HTMX + debounced POST to a `/preview` endpoint is idiomatic for this Go stack. |
| Create / edit / delete posts | Core CRUD — without this the admin panel has no purpose | Low | Standard forms + DB operations. Soft-delete preferred over hard-delete (recoverable). |
| Draft vs published state | Writing a post should not immediately make it public | Low | Boolean `published` field + `published_at` timestamp. Admin sees all; public sees only published. |

---

## Differentiators

Features that distinguish this blog. Not universally expected, but add real value when done well.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Thumbs-up / reaction counter | Low-friction reader engagement — no account required, no comment thread to moderate | Low | Single integer counter per post. Endpoint `POST /posts/:slug/react` increments in DB. Rate-limit by IP or cookie to prevent trivial gaming. No auth needed for readers. |
| API endpoint for pushing .md files | Enables local writing workflow (VS Code, Obsidian, Neovim) with `curl` or a script to publish without touching the browser | Medium | `POST /api/posts` with bearer token auth + multipart or JSON body containing frontmatter + markdown. Token stored in env var. Separate from session auth. |
| Table of contents auto-generation | On long technical posts, a TOC is a significant UX improvement; readers can jump to sections | Low | goldmark's `goldmark-toc` extension or custom AST walker to collect headings. Render as sidebar or inline at post top. |
| Estimated reading time | Small signal that sets reader expectations and reduces abandonment on long posts; Yoast, Ghost, Medium all include this | Low | Word count / 200 wpm. Calculate at render time or store on post. Display as "~5 min read". |
| LaTeX / math rendering | Differentiator for technical/academic content; most personal blogs don't bother | High | KaTeX is faster than MathJax (1.8s vs 4.2s for 500 equations) and sufficient for personal blog usage. Server-side rendering via Go → subprocess is complex; client-side KaTeX JS loading only when post has math is simpler and correct. Add a `has_math` boolean front-matter field to avoid loading KaTeX on every page. |
| Dark mode | ~50% of developer-audience readers prefer dark mode; CSS `prefers-color-scheme` media query costs almost nothing | Low | CSS variables + `@media (prefers-color-scheme: dark)`. No JS toggle needed for basic support; optional manual toggle as enhancement. |
| Post cover images | Makes listing pages visually distinct; improves social share unfurls | Low | Optional `cover_image` field per post. Falls back to site-default OG image. |
| Sitemap (`/sitemap.xml`) | Helps search engines discover all posts; especially matters for new or low-authority domains | Low | Static generation on publish, or dynamic generation from DB query. Update whenever post is published/deleted. |
| Canonical URL tag | Prevents SEO dilution if content is syndicated elsewhere (e.g. dev.to cross-posts) | Low | `<link rel="canonical">` in `<head>`. Defaults to current URL; override field on post for syndication. |
| Post series / related posts | For multi-part technical deep dives, linking next/previous keeps readers in the series | Medium | Tags or explicit `series` field. Low-complexity if deferred to Phase 2+. |

---

## Anti-Features

Things to explicitly not build for this project. Each has a reason and an alternative.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| Comments system | Requires spam moderation, storage schema, reader identity management, and ongoing maintenance. For a personal blog the author is the only voice that needs to be heard. | Use the thumbs-up reaction counter for lightweight reader signal. Add a "contact me" link if readers need to reach out. |
| User registration / reader accounts | Zero value for a single-author personal blog. Adds auth complexity, GDPR surface, and password reset flows for users who gain nothing by having accounts. | Keep it public-read, session-only for admin. |
| OAuth / social login | Introduces third-party dependencies, token refresh complexity, and upstream breakage risk for a system with exactly one privileged user. | bcrypt + session cookie. Sufficient, boring, correct. |
| Tag/category filtering on public site (v1) | Without many posts, a tag cloud is noise. Premature taxonomy adds DB joins and template complexity without reader value. | Add a `tags` field to the schema from the start, but defer rendering tag listing pages until the post volume justifies them. |
| Full-text search | Postgres `pg_trgm` or external search (Meilisearch, Typesense) is non-trivial to operate on a single EBS instance. A blog with 50 posts doesn't need search — pagination + tag browsing is sufficient. | Browser `Ctrl+F` and Google site-search (`site:jared-wallace.com`) cover the actual need. Revisit when post count exceeds ~200. |
| Newsletter / email delivery | Requires an email service provider (Postmark, SES), subscription management, unsubscribe flows, CAN-SPAM compliance. Significant operational overhead. | RSS covers the follow-me use case. Add newsletter if audience demand becomes clear. |
| Analytics dashboard in-app | Building analytics is a project in itself. Maintaining it diverts effort from the actual blog. | Use Plausible, Fathom, or Cloudflare analytics (privacy-respecting, zero maintenance). |
| Image optimization / CDN pipeline | Sharp/libvips integration, WebP conversion, responsive srcsets — legitimate features for high-traffic sites, overkill for a personal blog on EBS. | Serve images directly from EBS via Go static handler. Revisit if image-heavy posts cause load problems. |
| WYSIWYG rich text editor (ProseMirror, TipTap, Quill) | Heavy JS dependency, complex serialization to/from markdown, diverges from the "write markdown" philosophy. The dual workflow (browser editor + API push) already covers all writing modes. | Stick to the split-pane raw markdown + preview editor. |
| Multi-author support | This is a personal blog. Multi-author adds roles, per-author profiles, attribution complexity, and admin privilege scoping for no benefit. | If a second author ever appears, revisit schema at that time. |

---

## Feature Dependencies

```
Markdown rendering
  └── Syntax highlighting (requires markdown pipeline to be in place)
  └── Table of contents (requires heading AST walk, built on renderer)
  └── LaTeX rendering (requires post-processing hook or front-matter flag)

Post CRUD
  └── Draft/published state (requires published_at + boolean in post schema)
  └── Slug management (requires post creation flow)
  └── Cover image field (requires post schema + static file serving)
  └── Tags field (requires post schema; display deferred)

Session auth
  └── Admin CRUD (auth must be in place before admin routes are exposed)
  └── Web editor (depends on admin routing)

API push endpoint
  └── Post CRUD (API is an alternative ingestion path for the same schema)
  └── Bearer token (separate from session auth, simpler to add after session auth)

RSS feed
  └── Post listing + published state (feed only includes published posts)
  └── Slug (required for item <link>)

Open Graph tags
  └── Cover image field (optional enhancement to OG image)
  └── Post published state (only published posts should be indexed)

Sitemap
  └── Post listing + slugs (generates from same data as listing)
  └── Published state (only published posts in sitemap)

Thumbs-up counter
  └── Post CRUD (counter lives on the post record or a separate reactions table)

LaTeX rendering
  └── Markdown rendering (LaTeX delimiters live inside markdown source)
  └── has_math front-matter field (to avoid loading KaTeX universally)
```

---

## MVP Recommendation

The required MVP per PROJECT.md is already well-scoped. Priority order based on dependency graph and reader value:

**Phase 1 — Core reading experience:**
1. Markdown rendering with goldmark (table stakes; everything else depends on it)
2. Syntax-highlighted code blocks (expected by technical audience from day one)
3. Post listing with pagination + published/draft state
4. Post published date and readable slugs
5. Responsive layout + 404 page

**Phase 2 — Admin and writing:**
6. Session auth (prerequisite for all admin features)
7. Post CRUD (create, edit, soft-delete)
8. Draft/published workflow
9. Web editor with split-pane markdown preview

**Phase 3 — Distribution and engagement:**
10. RSS feed (demanded by feed-reader users; simple to build once posts exist)
11. Open Graph meta tags + canonical URL
12. Thumbs-up reaction counter
13. API push endpoint for local `.md` files

**Phase 4 — Enhancements (stretch / post-v1):**
14. Sitemap
15. Estimated reading time
16. Table of contents
17. Dark mode (CSS-only)
18. Post cover images
19. LaTeX / KaTeX rendering (front-matter-gated)

**Defer indefinitely:**
- Comments, user accounts, OAuth, full-text search, newsletter, in-app analytics

---

## Sources

- [goldmark — CommonMark-compliant Go markdown parser](https://github.com/yuin/goldmark) — HIGH confidence
- [goldmark-highlighting/v2 — Chroma syntax highlighting extension](https://pkg.go.dev/github.com/yuin/goldmark-highlighting/v2) — HIGH confidence
- [How I developed a markdown blog in Go and HTMX (fluxsec.red)](https://fluxsec.red/how-I-developed-a-markdown-blog-with-go-and-HTMX) — MEDIUM confidence
- [Building a Markdown Blog Engine with Go and Fiber (dasroot.net, 2026)](https://dasroot.net/posts/2026/01/building-markdown-blog-engine-go-fiber/) — MEDIUM confidence
- [KaTeX vs MathJax performance (KaTeX docs)](https://katex.org/) — HIGH confidence
- [RSS 2.0 specification](https://validator.w3.org/feed/docs/rss2.html) — HIGH confidence
- [Open Graph protocol and SEO (bigredseo.com, 2025)](https://www.bigredseo.com/wordpress-seo-open-graph-setup/) — MEDIUM confidence
- [Yoast estimated reading time feature](https://yoast.com/features/estimated-reading-time/) — MEDIUM confidence
- [Mataroa — minimalist blogging reference](https://mataroa.blog/) — MEDIUM confidence
- [Why RSS feeds are still relevant in 2025](https://rss.app/blog/why-rss-feeds-are-still-relevant-in-2025-4h4WiW) — MEDIUM confidence
