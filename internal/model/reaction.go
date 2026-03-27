package model

import "time"

// Reaction represents a reader's thumbs-up on a blog post.
type Reaction struct {
	CreatedAt time.Time
	IPHash    string
	ID        int64
	PostID    int64
}
