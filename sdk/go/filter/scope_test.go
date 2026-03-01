package filter

import "testing"

func TestMatchScope(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		scope   string
		want    bool
	}{
		// Exact match (no wildcards)
		{"exact match", "cli", "cli", true},
		{"exact no match", "cli", "web", false},
		{"exact empty", "", "", true},
		{"exact empty pattern", "", "cli", false},

		// Star wildcard
		{"prefix wildcard", "cli*", "cli", true},
		{"prefix wildcard match", "cli*", "cli-tools", true},
		{"prefix wildcard no match", "cli*", "web", false},
		{"suffix wildcard", "*cli", "cli", true},
		{"suffix wildcard match", "*cli", "my-cli", true},
		{"suffix wildcard no match", "*cli", "web", false},
		{"contains wildcard", "*web*", "web", true},
		{"contains wildcard match", "*web*", "my-web-app", true},
		{"contains wildcard no match", "*web*", "cli", false},

		// Star wildcard crosses path separators
		{"prefix wildcard crosses slash", "cli*", "cli/import", true},
		{"prefix wildcard crosses slash nested", "cli*", "cli/tracks/foo", true},
		{"mid wildcard crosses slash", "*import*", "cli/import", true},
		{"suffix wildcard crosses slash", "*import", "cli/import", true},

		// Mixed literal slash and wildcard in pattern
		{"slash then wildcard", "cli/*", "cli/import", true},
		{"slash then wildcard nested no match", "cli/*", "web/import", false},
		{"partial after slash", "cli/imp*", "cli/import", true},
		{"partial after slash no match", "cli/imp*", "cli/export", false},

		// Multiple wildcards
		{"double wildcard", "*cli*import*", "my-cli/import-tool", true},
		{"double wildcard no match", "*cli*import*", "my-web/export", false},

		// Question mark and brackets treated as literal (not wildcards)
		{"question mark literal", "cl?", "cl?", true},
		{"question mark literal no match", "cl?", "cli", false},
		{"bracket literal", "cl[ij]", "cl[ij]", true},
		{"bracket literal no match", "cl[ij]", "cli", false},

		// Backward compat: no wildcard uses exact equality
		{"no wildcard exact", "cli/graph", "cli/graph", true},
		{"no wildcard no match", "cli/graph", "cli/next", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchScope(tt.pattern, tt.scope)
			if got != tt.want {
				t.Errorf("MatchScope(%q, %q) = %v, want %v", tt.pattern, tt.scope, got, tt.want)
			}
		})
	}
}

