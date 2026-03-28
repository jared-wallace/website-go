package post

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// slugRegex matches any sequence of characters that are not lowercase letters
// or digits. Compiled once at package init to avoid repeated allocation.
var slugRegex = regexp.MustCompile(`[^a-z0-9]+`)

// GenerateSlug converts a post title to a URL-friendly slug.
// Rules:
//  1. Lowercase the input.
//  2. Replace any run of non-alphanumeric characters with a single hyphen.
//  3. Trim leading/trailing hyphens.
//  4. If the result is empty, return "post-{unix_timestamp}" as a fallback.
func GenerateSlug(title string) string {
	s := strings.ToLower(title)
	s = slugRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return fmt.Sprintf("post-%d", time.Now().Unix())
	}
	return s
}
