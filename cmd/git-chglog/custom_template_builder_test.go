package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomTemplateBuilderDefault(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleGitHub,
		CommitMessageFormat: fmtTypeScopeSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Pull Requests

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderNone(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeScopeSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderSubjectOnly(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
{{ range .Commits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderSubject(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderIgnoreReverts(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       true,
		IncludeReverts:      false,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ .Subject }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderIgnoreMerges(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeSubject.display,
		Template:            tplStandard.display,
		IncludeMerges:       false,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestCustomTemplateBuilderCool(t *testing.T) {
	assert := assert.New(t)
	builder := NewCustomTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeScopeSubject.display,
		Template:            tplCool.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ range .Versions }}
## {{ .Tag.Name }}

> {{ datetime "2006-01-02" .Tag.Date }}

{{ range .CommitGroups -}}
### {{ .Title }}

{{ range .Commits -}}
* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges

{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}
