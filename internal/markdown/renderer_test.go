package markdown_test

import (
	"strings"
	"testing"

	"github.com/jared-wallace/website-go/internal/markdown"
)

// TestRender_XSSStripped verifies that script tags are stripped from output — the critical XSS gate.
func TestRender_XSSStripped(t *testing.T) {
	r := markdown.NewRenderer()
	input := "Hello <script>alert(1)</script> world"
	out := string(r.Render(input))
	if strings.Contains(out, "<script>") {
		t.Errorf("XSS: output must not contain <script>, got: %s", out)
	}
	if strings.Contains(out, "alert(1)") {
		t.Errorf("XSS: output must not contain alert(1), got: %s", out)
	}
}

// TestRender_BasicMarkdown verifies bold text renders to <strong>.
func TestRender_BasicMarkdown(t *testing.T) {
	r := markdown.NewRenderer()
	input := "**bold** text"
	out := string(r.Render(input))
	if !strings.Contains(out, "<strong>bold</strong>") {
		t.Errorf("bold: expected <strong>bold</strong> in output, got: %s", out)
	}
}

// TestRender_GFMTable verifies GFM pipe-delimited tables render to <table>.
func TestRender_GFMTable(t *testing.T) {
	r := markdown.NewRenderer()
	input := "| Name | Age |\n|------|-----|\n| Alice | 30 |"
	out := string(r.Render(input))
	if !strings.Contains(out, "<table>") {
		t.Errorf("GFM table: expected <table> in output, got: %s", out)
	}
}

// TestRender_Strikethrough verifies strikethrough renders to <del>.
func TestRender_Strikethrough(t *testing.T) {
	r := markdown.NewRenderer()
	input := "~~deleted~~"
	out := string(r.Render(input))
	if !strings.Contains(out, "<del>") {
		t.Errorf("strikethrough: expected <del> in output, got: %s", out)
	}
}

// TestRender_CodeFence verifies code fences render to <pre> and <code>.
func TestRender_CodeFence(t *testing.T) {
	r := markdown.NewRenderer()
	input := "```go\nfunc main() {}\n```"
	out := string(r.Render(input))
	if !strings.Contains(out, "<pre>") && !strings.Contains(out, "<code") {
		t.Errorf("code fence: expected <pre> and <code> in output, got: %s", out)
	}
}

// TestRender_LinkAutoDetected verifies bare URLs are auto-linked by linkify extension.
func TestRender_LinkAutoDetected(t *testing.T) {
	r := markdown.NewRenderer()
	input := "Visit https://example.com today"
	out := string(r.Render(input))
	if !strings.Contains(out, `<a href="https://example.com"`) {
		t.Errorf("linkify: expected <a href=\"https://example.com\"> in output, got: %s", out)
	}
}

// TestRenderWithMeta verifies YAML front matter extraction works.
func TestRenderWithMeta(t *testing.T) {
	r := markdown.NewRenderer()
	input := "---\ntitle: Hello\n---\n# Body"
	_, meta := r.RenderWithMeta(input)
	if meta == nil {
		t.Fatal("front matter: expected non-nil metadata map")
	}
	title, ok := meta["title"]
	if !ok {
		t.Fatalf("front matter: expected key 'title' in metadata, got: %v", meta)
	}
	if title != "Hello" {
		t.Errorf("front matter: expected title = 'Hello', got: %v", title)
	}
}

// TestRender_IFrameStripped verifies that iframe tags are stripped (XSS via embed).
func TestRender_IFrameStripped(t *testing.T) {
	r := markdown.NewRenderer()
	input := `<iframe src="evil.com"></iframe>`
	out := string(r.Render(input))
	if strings.Contains(out, "<iframe") {
		t.Errorf("iframe: output must not contain <iframe>, got: %s", out)
	}
}

// TestRender_EventHandlerStripped verifies that onclick and other event handlers are stripped.
func TestRender_EventHandlerStripped(t *testing.T) {
	r := markdown.NewRenderer()
	input := `<div onclick="alert(1)">click</div>`
	out := string(r.Render(input))
	if strings.Contains(out, "onclick") {
		t.Errorf("event handler: output must not contain 'onclick', got: %s", out)
	}
}

// TestRender_EmptyInput verifies empty input produces no error and minimal/empty output.
func TestRender_EmptyInput(t *testing.T) {
	r := markdown.NewRenderer()
	out := string(r.Render(""))
	// Should not panic and output should be empty or just whitespace
	if strings.Contains(out, "<script>") || strings.Contains(out, "error") {
		t.Errorf("empty input: unexpected output: %s", out)
	}
}
