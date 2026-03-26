package post

import (
	"strings"

	"golang.org/x/net/html"
)

// ToCEntry represents a single heading in a post's table of contents.
type ToCEntry struct {
	ID    string
	Text  string
	Level int // 2 or 3 (h2 or h3)
}

// tocMinEntries is the minimum heading count required to emit a ToC.
// Fewer than 3 headings is not worth displaying.
const tocMinEntries = 3

// ExtractToC walks the rendered HTML and collects all h2/h3 elements that
// carry an id attribute. Returns nil when fewer than tocMinEntries headings
// are found — callers can treat nil as "no ToC to render".
func ExtractToC(renderedHTML string) []ToCEntry {
	root, err := html.Parse(strings.NewReader(renderedHTML))
	if err != nil {
		return nil
	}

	var entries []ToCEntry
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			level := headingLevel(n.Data)
			if level == 2 || level == 3 {
				id := attrVal(n, "id")
				if id != "" {
					entries = append(entries, ToCEntry{
						ID:    id,
						Text:  textContent(n),
						Level: level,
					})
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(root)

	if len(entries) < tocMinEntries {
		return nil
	}
	return entries
}

// headingLevel maps "h2" → 2, "h3" → 3, other → 0.
func headingLevel(tag string) int {
	switch tag {
	case "h2":
		return 2
	case "h3":
		return 3
	default:
		return 0
	}
}

// attrVal returns the value of the named attribute, or "" if absent.
func attrVal(n *html.Node, name string) string {
	for _, a := range n.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// textContent returns the concatenated text content of all descendant text nodes.
func textContent(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return sb.String()
}
