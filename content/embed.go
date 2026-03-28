// Package content holds embedded static content files for the site.
// Files are baked into the binary at build time via go:embed.
package content

import _ "embed"

// AboutMD is the raw Markdown source for the /about page.
// Rendered at request time through the goldmark pipeline.
//
//go:embed about.md
var AboutMD string
