package slug

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic transformations
		{"lowercase", "Hello World", "hello-world"},
		{"spaces to hyphens", "foo bar baz", "foo-bar-baz"},
		{"already lowercase", "hello", "hello"},

		// Special characters
		{"at sign", "user@domain", "user-domain"},
		{"exclamation", "hello!", "hello"},
		{"slash", "path/to/thing", "path-to-thing"},
		{"underscore", "snake_case", "snake-case"},
		{"mixed special chars", "a@b!c/d_e", "a-b-c-d-e"},

		// Trimming leading/trailing hyphens
		{"leading special", "@hello", "hello"},
		{"trailing special", "hello!", "hello"},
		{"both sides special", "!hello!", "hello"},

		// Consecutive special chars collapse to single hyphen
		{"consecutive specials", "a---b", "a-b"},
		{"multiple spaces", "a   b", "a-b"},

		// Length truncation at 50 chars
		{"exactly 50", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghij", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghij"},
		{"over 50 truncated", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijXXX", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghij"},
		{"trailing hyphen after truncation", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghi stuff", "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghi"},

		// Edge cases
		{"empty string", "", ""},
		{"all special chars", "@!#$%", ""},
		{"single char", "A", "a"},
		{"numbers", "task123", "task123"},
		{"numbers only", "12345", "12345"},
		{"mixed case with numbers", "Task 42 Done", "task-42-done"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
