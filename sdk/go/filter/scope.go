package filter

import "strings"

// MatchScope checks whether a scope value matches a pattern.
// The only supported wildcard is '*', which matches zero or more
// characters including '/'. Scopes are logical identifiers
// (e.g. "cli/import") so "cli*" matches "cli", "cli/import", etc.
// Without a wildcard, exact string equality is used.
func MatchScope(pattern, scope string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == scope
	}
	return matchStar(pattern, scope)
}

// matchStar matches a pattern containing '*' wildcards against s.
// Each '*' matches zero or more arbitrary characters (including '/').
func matchStar(pattern, s string) bool {
	parts := strings.Split(pattern, "*")
	// Fast path: single '*' prefix/suffix/contains.
	if len(parts) == 2 {
		return strings.HasPrefix(s, parts[0]) && strings.HasSuffix(s, parts[1]) && len(s) >= len(parts[0])+len(parts[1])
	}
	// General case: parts must appear in order within s.
	// First part must be a prefix, last part must be a suffix.
	if !strings.HasPrefix(s, parts[0]) {
		return false
	}
	s = s[len(parts[0]):]
	for _, part := range parts[1 : len(parts)-1] {
		idx := strings.Index(s, part)
		if idx < 0 {
			return false
		}
		s = s[idx+len(part):]
	}
	return strings.HasSuffix(s, parts[len(parts)-1])
}
