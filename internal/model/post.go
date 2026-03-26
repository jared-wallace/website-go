package model

import "time"

// Post represents a blog post. It maps 1:1 to the posts table in Postgres.
// The RenderedHTML field is pre-computed at write time from Body (markdown)
// and stored to avoid re-rendering on every read.
type Post struct {
	ID           int64
	Title        string
	Slug         string
	Body         string
	RenderedHTML string
	Published    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time // soft delete (ADMIN-05)
}
