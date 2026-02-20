package cli

import (
	_ "embed"

	"github.com/driangle/taskmd/apps/cli/internal/template"
)

//go:embed templates/task_feature.md
var taskFeatureTemplate string

//go:embed templates/task_bug.md
var taskBugTemplate string

//go:embed templates/task_chore.md
var taskChoreTemplate string

func init() {
	template.BuiltinTemplates = map[string]string{
		"feature": taskFeatureTemplate,
		"bug":     taskBugTemplate,
		"chore":   taskChoreTemplate,
	}
}
