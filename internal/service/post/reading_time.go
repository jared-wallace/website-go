package post

import (
	"math"
	"strings"
)

const wordsPerMinute = 200

// ReadingTime returns the estimated reading time in minutes for the given
// markdown body. It operates on the raw markdown text (not rendered HTML) so
// the word count reflects the actual prose, not tag noise.
// Minimum return value is 1 (a post with zero words still takes a minute).
func ReadingTime(body string) int {
	words := len(strings.Fields(body))
	if words == 0 {
		return 1
	}
	minutes := int(math.Ceil(float64(words) / float64(wordsPerMinute)))
	if minutes < 1 {
		return 1
	}
	return minutes
}
