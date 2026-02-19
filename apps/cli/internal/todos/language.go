package todos

// CommentSyntax describes how comments are written in a language.
type CommentSyntax struct {
	LinePrefix []string // e.g. ["//"], ["#"]
	BlockStart string   // e.g. "/*"
	BlockEnd   string   // e.g. "*/"
}

var languageMap = map[string]*CommentSyntax{
	// C-style languages
	".go":   {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".js":   {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".jsx":  {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".ts":   {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".tsx":  {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".rs":   {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},
	".css":  {BlockStart: "/*", BlockEnd: "*/"},
	".scss": {LinePrefix: []string{"//"}, BlockStart: "/*", BlockEnd: "*/"},

	// Hash-style languages
	".py":   {LinePrefix: []string{"#"}, BlockStart: `"""`, BlockEnd: `"""`},
	".rb":   {LinePrefix: []string{"#"}, BlockStart: "=begin", BlockEnd: "=end"},
	".sh":   {LinePrefix: []string{"#"}},
	".bash": {LinePrefix: []string{"#"}},
	".zsh":  {LinePrefix: []string{"#"}},
	".yaml": {LinePrefix: []string{"#"}},
	".yml":  {LinePrefix: []string{"#"}},
	".toml": {LinePrefix: []string{"#"}},

	// Markup
	".html": {BlockStart: "<!--", BlockEnd: "-->"},
	".xml":  {BlockStart: "<!--", BlockEnd: "-->"},
}

// LookupSyntax returns the comment syntax for a file extension, or nil if unsupported.
func LookupSyntax(ext string) *CommentSyntax {
	return languageMap[ext]
}

var languageNameMap = map[string]string{
	".go":   "go",
	".js":   "javascript",
	".jsx":  "javascript",
	".ts":   "typescript",
	".tsx":  "typescript",
	".rs":   "rust",
	".css":  "css",
	".scss": "scss",
	".py":   "python",
	".rb":   "ruby",
	".sh":   "shell",
	".bash": "shell",
	".zsh":  "shell",
	".yaml": "yaml",
	".yml":  "yaml",
	".toml": "toml",
	".html": "html",
	".xml":  "xml",
}

// LookupLanguage returns the human-readable language name for a file extension.
func LookupLanguage(ext string) string {
	return languageNameMap[ext]
}
