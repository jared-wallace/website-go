---
phase: quick
plan: 260327-usk
type: execute
wave: 1
depends_on: []
files_modified: [README.md]
autonomous: true
requirements: []
must_haves:
  truths:
    - "A developer visiting the repo sees a clear, accurate README describing the project"
    - "README covers what the project is, how to run it locally, and how to develop"
  artifacts:
    - path: "README.md"
      provides: "Project documentation"
      min_lines: 40
  key_links: []
---

<objective>
Create a README.md for the website-go repository.

Purpose: The repo has no README — visitors (and future-you) get zero orientation. A good README answers "what is this, how do I run it, and how do I hack on it" in under 2 minutes of reading.
Output: README.md at repo root
</objective>

<execution_context>
@$HOME/.claude/get-shit-done/workflows/execute-plan.md
@$HOME/.claude/get-shit-done/templates/summary.md
</execution_context>

<context>
@CLAUDE.md
@go.mod
@Makefile
@docker-compose.dev.yml
</context>

<tasks>

<task type="auto">
  <name>Task 1: Create README.md</name>
  <files>README.md</files>
  <action>
Create README.md at the repo root with the following sections. Keep it concise — a weathered beach bar doesn't need a 200-line manifesto.

**Structure:**

1. **Title + one-liner** — `# website-go` followed by a single sentence: personal blog platform for jared-wallace.com, Go web server with a "weathered beach bar" nautical theme.

2. **Tech Stack** — Brief bullet list of key technologies (Go 1.26, PostgreSQL, goldmark, pgx/v5, goose migrations, html/template, scs sessions). No version numbers needed here — go.mod is the source of truth.

3. **Prerequisites** — Go 1.23+ (local dev) or Go 1.26 (Docker build), PostgreSQL 15+, Docker and docker-compose (for local Postgres). Mention `golangci-lint` for linting.

4. **Getting Started** — Step-by-step:
   - Clone the repo
   - `make dev-up` to start Postgres
   - Set `DATABASE_URL` env var (provide the default from docker-compose.dev.yml: `postgres://website:website@localhost:5432/website?sslmode=disable`)
   - `make migrate` to run migrations
   - `make run` to build and start the server
   - Visit `http://localhost:8080`

5. **Development** — Reference Makefile targets with brief descriptions. Mention `make help` for the full list. Note that `make dev` uses `air` for hot reload.

6. **Project Structure** — Brief directory layout:
   ```
   cmd/server/     — Application entrypoint and wiring
   internal/       — Core packages (config, database, handler, markdown, model, server)
   db/migrations/  — Goose SQL migrations (embedded at build time)
   web/            — Templates and static assets
   ```

7. **License** — "All rights reserved" or omit if no license file exists.

Do NOT include: badges, CI status, contributor guidelines, code of conduct, or deployment instructions (that's Phase 6 territory).

Run the result through a simplicity check — every sentence should earn its keep.
  </action>
  <verify>
    <automated>test -f README.md && wc -l README.md | awk '{if ($1 >= 40) print "OK"; else print "FAIL: too short"}'</automated>
  </verify>
  <done>README.md exists at repo root, covers project description, tech stack, prerequisites, getting started, development workflow, and project structure. Reads clearly in under 2 minutes.</done>
</task>

</tasks>

<verification>
- README.md exists and renders valid Markdown
- All referenced Makefile targets (`dev-up`, `migrate`, `run`, `dev`, `help`) actually exist in Makefile
- No broken links or references to nonexistent files
</verification>

<success_criteria>
README.md exists, is accurate to the current codebase, and a new developer could go from clone to running server by following it.
</success_criteria>

<output>
After completion, create `.planning/quick/260327-usk-update-readme-md/260327-usk-SUMMARY.md`
</output>
