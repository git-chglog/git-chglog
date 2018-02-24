package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplateBuilderDefault(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleGitHub,
		CommitMessageFormat: fmtTypeScopeSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
<a name="{{.Tag.Name}}"></a>
## {{if .Tag.Previous}}[{{.Tag.Name}}]({{$.Info.RepositoryURL}}/compare/{{.Tag.Previous.Name}}...{{.Tag.Name}}){{else}}{{.Tag.Name}}{{end}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{if ne .Scope ""}}**{{.Scope}}:** {{end}}{{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Pull Requests
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderNone(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeScopeSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{if ne .Scope ""}}**{{.Scope}}:** {{end}}{{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Merges
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderCool(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeScopeSubject.Display,
		Template:            tplCool,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}}

> {{datetime "2006-01-02" .Tag.Date}}
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{if ne .Scope ""}}**{{.Scope}}:** {{end}}{{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Merges
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderSubjectOnly(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
{{range .Commits}}
* {{.Header}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Merges
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderSubject(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{if .MergeCommits}}
### Merges
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderIgnoreReverts(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       true,
		IncludeReverts:      false,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{.Subject}}{{end}}
{{end}}{{if .MergeCommits}}
### Merges
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}

func TestTemplateBuilderIgnoreMerges(t *testing.T) {
	assert := assert.New(t)
	builder := NewTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.Display,
		Template:            tplStandard,
		IncludeMerges:       false,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{range .Versions}}
## {{.Tag.Name}} ({{datetime "2006-01-02" .Tag.Date}})
{{range .CommitGroups}}
### {{.Title}}
{{range .Commits}}
* {{.Subject}}{{end}}
{{end}}{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`, out)
}
