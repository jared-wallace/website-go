// Package migrations embeds the SQL migration files for goose.
// Placing the embed here (alongside the SQL files) avoids the go:embed
// restriction against ".." path components.
package migrations

import "embed"

// FS holds all goose SQL migration files embedded at build time.
//
//go:embed *.sql
var FS embed.FS
