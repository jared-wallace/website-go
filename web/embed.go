package web

import "embed"

// Templates holds the embedded HTML template files under web/templates/.
// Loaded at build time — no runtime file I/O required.
//
//go:embed templates
var Templates embed.FS

// Static holds the embedded static asset files under web/static/.
// Loaded at build time — no runtime file I/O required.
//
//go:embed static
var Static embed.FS
