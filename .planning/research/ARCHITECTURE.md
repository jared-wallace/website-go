# Architecture Patterns

**Domain:** Personal blog platform (Go, server-side rendered, single-admin)
**Researched:** 2026-03-26
**Overall confidence:** HIGH — patterns are well-established; specific library versions verified against official sources

---

## Recommended Architecture

A three-layer monolith served from a single Go binary, behind Nginx + AWS ALB. The binary handles both the public blog (jared-wallace.com) and the admin panel (admin.jared-wallace.com) via host-based routing inside the same process. One docker-compose stack: the Go app + a Postgres sidecar.

```
[Browser]
    |
    v
[AWS ALB] -- TLS termination (ACM cert)
    |
    v
[Nginx on EC2 :443/:80] -- reverse proxy
    |
    v
[Go app on :8080]
    |
    |-- /              --> Public blog handlers
    |-- /admin/...     --> Admin handlers (session-gated)
    |-- /api/push      --> API push handler (token-gated)
    |-- /feed.xml      --> RSS handler
    |-- /uploads/...   --> Static file server (EBS-mounted images)
    |-- /static/...    --> Embedded static assets (CSS, JS, fonts)
    |
    v
[Postgres sidecar :5432]   (data on EBS bind-mount at /var/www/html/pgdata)
```

---

## Component Boundaries

| Component | Responsibility | Communicates With |
|-----------|---------------|-------------------|
| `cmd/server/main.go` | Process entry point, wires dependencies, starts HTTP server | All internal packages |
| `internal/config` | Loads config from env vars or file at startup | `main.go` only |
| `internal/router` | Registers all routes, applies middleware chains | Handlers, middleware |
| `internal/middleware` | Auth guard, session loading, request logger, recovery | Router, sessions |
| `internal/handler/blog` | Public post list, single post, thumbs-up | PostService |
| `internal/handler/admin` | Login, editor, post CRUD, image upload | PostService, SessionStore |
| `internal/handler/api` | `POST /api/push` for local-to-server .md push | PostService |
| `internal/handler/feed` | Generates RSS XML | PostService |
| `internal/service/post` | Business logic: parse markdown, frontmatter, slugs | PostRepository |
| `internal/service/auth` | bcrypt verify, session create/destroy | SessionStore, DB |
| `internal/repository/post` | SQL queries: posts table CRUD | `*sql.DB` / sqlx |
| `internal/repository/session` | Session store backed by Postgres table | `*sql.DB` |
| `internal/markdown` | Goldmark conversion, sanitization, syntax highlight | PostService |
| `db/migrations/` | SQL migration files (embedded via `//go:embed`) | golang-migrate at startup |
| `web/templates/` | html/template files for blog and admin views | Handlers |
| `web/static/` | CSS, JS, fonts (embedded into binary via `//go:embed`) | Router |
| `docker-compose.yml` | Defines `app` + `db` services, bind-mounts EBS path | EC2 host |

---

## Data Flow

### Public: Reader Loads a Post

```
Browser GET /posts/my-slug
  -> Router -> blog.Handler.ShowPost
  -> PostService.GetBySlug(slug)
    -> PostRepository.FindBySlug(slug)   [SQL SELECT]
    <- Post{Body: "# Hello..."}
  -> markdown.Render(post.Body)          [Goldmark HTML]
  -> template.Execute(w, data)           [html/template]
Browser receives rendered HTML page
```

### Admin: Publish a Post via Web Editor

```
Admin POST /admin/posts (form submit)
  -> middleware.RequireSession            [reject if no valid session]
  -> admin.Handler.CreatePost
  -> PostService.Create(form data)
    -> markdown.ParseFrontmatter(body)   [extract title, date, tags]
    -> PostRepository.Insert(post)       [SQL INSERT]
  -> redirect to /admin/posts
```

### API: Push Markdown from Local Machine

```
curl -X POST /api/push -H "Authorization: Bearer $TOKEN" --data-binary @post.md
  -> middleware.RequireAPIToken           [static token from config]
  -> api.Handler.Push
  -> PostService.Upsert(markdown bytes)
    -> markdown.ParseFrontmatter(body)
    -> PostRepository.Upsert(post)       [INSERT ... ON CONFLICT DO UPDATE]
  -> 200 OK
```

### Session Auth Flow

```
POST /admin/login (username + password)
  -> auth.Service.Authenticate(username, password)
    -> DB lookup user row
    -> bcrypt.CompareHashAndPassword
  -> scs.SessionManager.RenewToken(ctx)  [session fixation prevention]
  -> scs.Put(ctx, "userID", id)
  -> redirect to /admin
```

---

## Patterns to Follow

### Pattern: Dependency Injection via Constructor (not global state)

Wire dependencies at startup in `main.go`. Pass them down. No `init()` globals.

```go
// main.go
db := mustConnect(cfg.DatabaseURL)
postRepo := repository.NewPostRepo(db)
postSvc  := service.NewPostService(postRepo, markdown.NewRenderer())
blogH    := handler.NewBlogHandler(postSvc)
adminH   := handler.NewAdminHandler(postSvc, sessionManager)
router   := router.New(blogH, adminH, sessionManager)
```

### Pattern: Thin Handlers, Fat Services

Handlers translate HTTP to/from domain types. Business logic lives in services. Repositories handle only SQL.

```go
// handler — HTTP boundary only
func (h *BlogHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
    slug := chi.URLParam(r, "slug")
    post, err := h.posts.GetBySlug(r.Context(), slug)
    if err != nil { /* handle */ return }
    h.tmpl.Render(w, "post.html", post)
}

// service — domain logic
func (s *PostService) GetBySlug(ctx context.Context, slug string) (*Post, error) {
    p, err := s.repo.FindBySlug(ctx, slug)
    if err != nil { return nil, err }
    p.RenderedHTML = s.renderer.Render(p.Body)
    return p, nil
}
```

### Pattern: Middleware Chain via chi (or net/http 1.22)

```
Request
  -> Logger (always)
  -> Recover (always)
  -> LoadSession (always, via scs middleware)
  -> [public routes — no auth]
  -> [/admin routes — RequireSession middleware]
  -> [/api routes  — RequireAPIToken middleware]
Handler
```

Go 1.22's enhanced `net/http.ServeMux` is sufficient for this project's routing needs (no complex regex patterns). Use chi if route grouping + per-group middleware becomes unwieldy — it is a non-painful refactor.

### Pattern: Templates with `//go:embed`

Embed `web/templates/` and `web/static/` into the binary at compile time. This keeps the Docker image self-contained with no volume mounts for static assets.

```go
//go:embed web/templates
var templateFS embed.FS

//go:embed web/static
var staticFS embed.FS
```

Uploaded images (user content) are NOT embedded — they live on the EBS bind-mount at `/var/www/html/uploads/` and are served via `http.FileServer`.

### Pattern: Migrations Embedded, Run at Startup (with guard)

Embed SQL migrations via `//go:embed db/migrations` and run via `golang-migrate` at startup before the server starts accepting traffic. For a single-instance deployment this is safe and removes operational complexity of a separate migration step.

---

## Anti-Patterns to Avoid

### Anti-Pattern: Global `db` or `cfg` Variables

**What:** Package-level `var DB *sql.DB` shared across packages.
**Why bad:** Hidden dependencies, untestable, race-prone.
**Instead:** Constructor injection — all database access flows through `*Repository` types initialized in `main.go`.

### Anti-Pattern: Rendering Markdown in the Handler

**What:** Calling `goldmark.Convert()` inside an HTTP handler on every request.
**Why bad:** Redundant CPU work; markdown output is deterministic for a given input.
**Instead:** Render markdown once on write (when post is saved) and store rendered HTML in the database. Serve the pre-rendered HTML on read. Regenerate only on update.

### Anti-Pattern: Storing Sessions in a Cookie Store (client-side)

**What:** Using gorilla/sessions with `CookieStore`, where all session data is encrypted into the cookie itself.
**Why bad:** Session data bloat; cannot invalidate sessions server-side; gorilla toolkit maintenance concerns.
**Instead:** Use `alexedwards/scs` with a Postgres-backed store (or in-memory store for single-instance). SCS is actively maintained, uses a server-side session table, and integrates cleanly as middleware.

### Anti-Pattern: Serving Uploads from a Non-Persistent Path

**What:** Writing uploaded images to the container filesystem at `/app/uploads/`.
**Why bad:** Docker container is ephemeral — files are lost on redeploy.
**Instead:** Bind-mount the EBS volume path into the container: `- /var/www/html/uploads:/app/uploads`. The Go server writes to `/app/uploads`; the host persists to EBS.

### Anti-Pattern: Running as Root Inside the Container

**What:** Default Dockerfile with no `USER` directive.
**Why bad:** Security risk; any container escape = root on host.
**Instead:** Add a non-root user in the Dockerfile. Ensure the bind-mount directory is chowned correctly on the host.

---

## Directory Layout

```
website-go/
├── cmd/
│   └── server/
│       └── main.go           # Entry point, DI wiring
├── internal/
│   ├── config/               # Env-var config loading
│   ├── handler/
│   │   ├── blog/             # Public read handlers
│   │   ├── admin/            # Admin CRUD handlers
│   │   ├── api/              # Push API handler
│   │   └── feed/             # RSS handler
│   ├── middleware/           # Session loader, auth guards, logger
│   ├── service/
│   │   ├── post/             # Post business logic
│   │   └── auth/             # Login, bcrypt
│   ├── repository/
│   │   ├── post/             # posts table queries (sqlx)
│   │   └── session/          # sessions table (scs pgxstore or sqlxstore)
│   └── markdown/             # Goldmark wrapper, sanitizer
├── db/
│   └── migrations/           # *.sql files (embedded)
├── web/
│   ├── templates/            # html/template files (embedded)
│   └── static/               # CSS, JS, fonts (embedded)
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

---

## Suggested Build Order (Phase Dependencies)

The components below must be built in dependency order — each layer depends on the one beneath it.

```
1. Config + DB connection + migrations
      |
2. Repository layer (SQL queries, testable with real Postgres)
      |
3. Markdown renderer (pure function, no deps)
      |
4. Service layer (depends on repo + renderer)
      |
5. Session management + auth service
      |
6. Middleware chain (depends on session manager)
      |
7. Public blog handlers + templates (depends on post service)
      |
8. Admin handlers + templates (depends on post service + auth)
      |
9. API push handler (depends on post service)
      |
10. RSS feed handler (depends on post service)
      |
11. Image upload + static file serving (depends on config for paths)
      |
12. Docker-compose + Dockerfile (wraps everything)
```

**Rationale for ordering:**
- Migrations must run before any repository code can execute against the schema.
- The repository interface should be defined alongside the service layer so the service can depend on an interface, not a concrete type — enabling test doubles.
- Auth and sessions come before admin handlers because the middleware that gates admin routes depends on them.
- Public handlers come before admin handlers — simpler, no auth dependency, validates the template/render pipeline first.
- Docker/deployment is last — it wraps a working binary, not a prototype.

---

## Scalability Considerations

This project targets a single-instance deployment. The architecture choices reflect that.

| Concern | At current scale (1 instance) | If ever scaling out |
|---------|-------------------------------|---------------------|
| Sessions | Postgres-backed scs store — works with 1 instance | Already shareable across instances (no sticky sessions needed) |
| Image uploads | EBS bind-mount on single EC2 | Would need S3 or shared NFS — known limitation, documented out of scope |
| DB connections | sqlx pool default (max 10–25) | Sufficient; PgBouncer if needed later |
| Caching | None needed — low traffic blog | Add Redis or in-process cache for rendered HTML if traffic warrants |
| Markdown rendering | Pre-render on write, serve stored HTML | Already efficient |

---

## Sources

- [Go Project Structure: Practices & Patterns (2025)](https://www.glukhov.org/post/2025/12/go-project-structure/)
- [The Repository Pattern in Go — Three Dots Labs](https://threedots.tech/post/repository-pattern-in-go/)
- [go-chi/chi — lightweight composable router](https://github.com/go-chi/chi)
- [Which Go Router Should I Use? — Alex Edwards](https://www.alexedwards.net/blog/which-go-router-should-i-use)
- [Go's 1.22+ ServeMux vs Chi Router — Calhoun.io](https://www.calhoun.io/go-servemux-vs-chi/)
- [alexedwards/scs — HTTP Session Management for Go](https://github.com/alexedwards/scs)
- [yuin/goldmark — CommonMark-compliant markdown parser](https://github.com/yuin/goldmark)
- [Go Database Patterns: GORM, sqlx, and pgx Compared](https://dasroot.net/posts/2025/12/go-database-patterns-gorm-sqlx-pgx-compared/)
- [golang-migrate — database migration tool](https://betterstack.com/community/guides/scaling-go/golang-migrate/)
- [Routing Enhancements for Go 1.22 — Official Go Blog](https://go.dev/blog/routing-enhancements)
- [Persisting container data — Docker Docs](https://docs.docker.com/get-started/docker-concepts/running-containers/persisting-container-data/)
