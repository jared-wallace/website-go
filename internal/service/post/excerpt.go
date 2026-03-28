package post

import (
	"regexp"
	"strings"
)

var (
	// Strip markdown link syntax: [text](url) → text
	reLink = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
	// Strip bold: **text** → text (must come before italic)
	reBold = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	// Strip italic: *text* → text
	reItalic = regexp.MustCompile(`\*([^*]+)\*`)
	// Strip inline code: `code` → code
	reCode = regexp.MustCompile("`([^`]+)`")
	// Strip ATX headings: ## heading → heading
	reHeading = regexp.MustCompile(`(?m)^#{1,6}\s+`)
	// Strip horizontal rules
	reHR = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)
	// Strip blockquote markers
	reBlockquote = regexp.MustCompile(`(?m)^>\s?`)
)

// Excerpt generates a plain-text preview from a markdown body, truncated to
// at most maxLen runes. If truncated, the result ends with "...". Markdown
// syntax (links, bold, italic, code fences, headings) is stripped so the
// excerpt reads as natural prose.
func Excerpt(markdownBody string, maxLen int) string {
	s := markdownBody

	// Apply markdown stripping in order
	s = reHR.ReplaceAllString(s, "")
	s = reHeading.ReplaceAllString(s, "")
	s = reBlockquote.ReplaceAllString(s, "")
	s = reLink.ReplaceAllString(s, "$1")
	s = reBold.ReplaceAllString(s, "$1")
	s = reItalic.ReplaceAllString(s, "$1")
	s = reCode.ReplaceAllString(s, "$1")

	// Normalize whitespace
	s = strings.Join(strings.Fields(s), " ")

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}

	// Truncate to maxLen then trim to last word boundary
	truncated := string(runes[:maxLen])
	if idx := strings.LastIndex(truncated, " "); idx > 0 {
		truncated = truncated[:idx]
	}
	return truncated + "..."
}

// ParseTags splits a comma-separated tag string into a slice of trimmed,
// non-empty tag values. Returns an empty (non-nil) slice for empty input.
func ParseTags(tagStr string) []string {
	parts := strings.Split(tagStr, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}
