package post

import (
	"strings"
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic title",
			input: "My First Post",
			want:  "my-first-post",
		},
		{
			name:  "hello world with punctuation",
			input: "Hello, World!",
			want:  "hello-world",
		},
		{
			name:  "spaces and dashes",
			input: "  spaces  and---dashes  ",
			want:  "spaces-and-dashes",
		},
		{
			name:  "uppercase",
			input: "UPPERCASE",
			want:  "uppercase",
		},
		{
			name:  "special characters",
			input: "special!@#$%chars",
			want:  "special-chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.input)
			if got != tt.want {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGenerateSlugEmpty(t *testing.T) {
	got := GenerateSlug("")
	if !strings.HasPrefix(got, "post-") {
		t.Errorf("GenerateSlug(\"\") = %q, want prefix \"post-\"", got)
	}
	// Should have a timestamp suffix that is numeric
	suffix := strings.TrimPrefix(got, "post-")
	if len(suffix) == 0 {
		t.Errorf("GenerateSlug(\"\") = %q, missing timestamp suffix", got)
	}
	for _, c := range suffix {
		if c < '0' || c > '9' {
			t.Errorf("GenerateSlug(\"\") = %q, suffix %q contains non-digit char %q", got, suffix, string(c))
		}
	}
}
