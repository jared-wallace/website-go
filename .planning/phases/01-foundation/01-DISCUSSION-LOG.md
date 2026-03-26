# Phase 1: Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-03-26
**Phase:** 01-foundation
**Areas discussed:** Project layout, Markdown pipeline, Makefile & dev workflow, GHA CI

---

## Project Layout

### Internal package organization

| Option | Description | Selected |
|--------|-------------|----------|
| Flat internal/ | internal/server/, internal/handler/, internal/markdown/, internal/database/, internal/model/ — simple, no nesting | :heavy_check_mark: |
| Domain-grouped internal/ | internal/blog/, internal/admin/, internal/auth/ — each with own handlers, models, queries | |
| You decide | Claude picks whichever structure fits best | |

**User's choice:** Flat internal/
**Notes:** None

### Asset embedding strategy

| Option | Description | Selected |
|--------|-------------|----------|
| go:embed | Templates and static files baked into the binary. Single artifact to deploy. | :heavy_check_mark: |
| Disk-based serving | Templates read from disk at runtime. Easy dev iteration. | |
| Hybrid | go:embed for prod, disk fallback in dev via build tag | |

**User's choice:** go:embed
**Notes:** None

### Configuration source

| Option | Description | Selected |
|--------|-------------|----------|
| Environment variables only | os.Getenv() with helpers. 12-factor, Docker-native. | :heavy_check_mark: |
| Config file + env override | YAML/TOML with env var precedence | |
| You decide | Claude picks simplest approach | |

**User's choice:** Environment variables only
**Notes:** None

---

## Markdown Pipeline

### Goldmark extensions

| Option | Description | Selected |
|--------|-------------|----------|
| GFM tables + strikethrough | Standard GitHub-flavored markdown | :heavy_check_mark: |
| Linkify | Auto-detect bare URLs | :heavy_check_mark: |
| Syntax highlighting (chroma) | Colored code fences | :heavy_check_mark: |
| YAML front matter (goldmark-meta) | Parse ---metadata--- blocks | :heavy_check_mark: |

**User's choice:** All four extensions
**Notes:** Multi-select — user enabled the full set

### bluemonday sanitization policy

| Option | Description | Selected |
|--------|-------------|----------|
| UGC policy | Allows common formatting, strips scripts/iframes/handlers | :heavy_check_mark: |
| Strict policy | Strips ALL HTML, plain text only | |
| Custom allowlist | UGC + specific extra tags (details, video, iframe) | |

**User's choice:** UGC policy
**Notes:** None

---

## Makefile & Dev Workflow

### Hot reload

| Option | Description | Selected |
|--------|-------------|----------|
| Yes, with air | cosmtrek/air watches files, rebuilds on change | :heavy_check_mark: |
| Manual restart | go run, rebuild manually | |
| You decide | Claude picks | |

**User's choice:** Yes, with air
**Notes:** None

### Local Postgres strategy

| Option | Description | Selected |
|--------|-------------|----------|
| docker-compose only | docker-compose.dev.yml with Postgres container | :heavy_check_mark: |
| System-installed Postgres | brew services, local install | |
| You decide | Claude picks | |

**User's choice:** docker-compose only
**Notes:** None

### Linter choice

| Option | Description | Selected |
|--------|-------------|----------|
| golangci-lint | Meta-linter with errcheck, govet, staticcheck, etc. | :heavy_check_mark: |
| Just go vet | Minimal, stdlib only | |
| You decide | Claude picks | |

**User's choice:** golangci-lint
**Notes:** None

---

## GHA CI

### CI database strategy

| Option | Description | Selected |
|--------|-------------|----------|
| Postgres service container | Real DB in CI, catches migration/query issues | :heavy_check_mark: |
| Unit tests only (no DB) | Faster, simpler, DB tests only local | |
| You decide | Claude picks | |

**User's choice:** Postgres service container
**Notes:** None

### CI job structure

| Option | Description | Selected |
|--------|-------------|----------|
| Single job | lint -> test -> build sequential. Conserves free-tier minutes. | :heavy_check_mark: |
| Parallel jobs | Separate lint, test, build jobs. Faster feedback, 3x minutes. | |
| You decide | Claude picks | |

**User's choice:** Single job
**Notes:** None

---

## Claude's Discretion

- Makefile target names and help text formatting
- .golangci.yml exact linter set beyond base six
- internal/config/ package shape
- Test file organization and naming
- goose migration numbering format

## Deferred Ideas

None — discussion stayed within phase scope.
