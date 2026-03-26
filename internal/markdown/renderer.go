package markdown

import (
	"bytes"
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// Renderer converts Markdown source to sanitized HTML using a goldmark pipeline
// followed by bluemonday sanitization. The pipeline order is critical:
// goldmark.Convert() THEN bluemonday.Sanitize(). Never reverse this.
type Renderer struct {
	gm goldmark.Markdown
	bm *bluemonday.Policy
}

// NewRenderer creates a Renderer with all four required extensions:
// GFM (tables + strikethrough), linkify, syntax highlighting, and YAML front matter.
func NewRenderer() *Renderer {
	gm := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,     // tables + strikethrough + autolinks
			extension.Linkify, // bare URL auto-linking (D-04)
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
			meta.Meta, // YAML front matter (D-04)
		),
		goldmark.WithRendererOptions(
			// html.WithUnsafe() allows raw HTML pass-through from goldmark.
			// This is safe ONLY because bluemonday runs downstream and strips
			// dangerous elements (scripts, iframes, event handlers, etc.).
			html.WithUnsafe(),
		),
	)

	// UGCPolicy (D-05): allows bold, italic, links, images, code, tables, blockquotes.
	// Strips <script>, <iframe>, <style>, event handlers, and javascript: URLs.
	bm := bluemonday.UGCPolicy()
	// Preserve class attributes on code/pre/span so syntax highlighting CSS classes survive.
	bm.AllowAttrs("class").OnElements("code", "span", "pre")

	return &Renderer{gm: gm, bm: bm}
}

// Render converts Markdown src to sanitized template.HTML.
// The returned value is safe to embed directly in Go HTML templates.
func (r *Renderer) Render(src string) template.HTML {
	var buf bytes.Buffer
	if err := r.gm.Convert([]byte(src), &buf); err != nil {
		return template.HTML("<!-- render error -->")
	}
	sanitized := r.bm.SanitizeBytes(buf.Bytes())
	return template.HTML(sanitized)
}

// RenderWithMeta converts Markdown src to sanitized template.HTML and also extracts
// YAML front matter into a metadata map. Returns nil metadata if none is present.
func (r *Renderer) RenderWithMeta(src string) (template.HTML, map[string]interface{}) {
	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := r.gm.Convert([]byte(src), &buf, parser.WithContext(ctx)); err != nil {
		return template.HTML("<!-- render error -->"), nil
	}
	metadata := meta.Get(ctx)
	sanitized := r.bm.SanitizeBytes(buf.Bytes())
	return template.HTML(sanitized), metadata
}
