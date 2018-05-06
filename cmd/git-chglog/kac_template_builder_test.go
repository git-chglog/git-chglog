package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKACTemplateBuilderDefault(t *testing.T) {
	assert := assert.New(t)
	builder := NewKACTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleGitHub,
		CommitMessageFormat: fmtTypeScopeSubject.display,
		Template:            tplKeepAChangelog.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ if .Versions -}}
<a name="unreleased"></a>
## [Unreleased]

{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Pull Requests
{{ range .MergeCommits -}}
- {{ .Header }}
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
{{ end -}}

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ end -}}
{{ end -}}`, out)
}

func TestKACTemplateBuilderNone(t *testing.T) {
	assert := assert.New(t)
	builder := NewKACTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtTypeScopeSubject.display,
		Template:            tplKeepAChangelog.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ if .Versions -}}
## Unreleased

{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ range .Versions }}
## {{ .Tag.Name }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges
{{ range .MergeCommits -}}
- {{ .Header }}
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

func TestKACTemplateBuilderSubject(t *testing.T) {
	assert := assert.New(t)
	builder := NewKACTemplateBuilder()

	out, err := builder.Build(&Answer{
		Style:               styleNone,
		CommitMessageFormat: fmtSubject.display,
		Template:            tplKeepAChangelog.display,
		IncludeMerges:       true,
		IncludeReverts:      true,
	})

	assert.Nil(err)
	assert.Equal(`{{ if .Versions -}}
## Unreleased

{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
{{ range .Commits -}}
- {{ .Header }}
{{ end }}
{{ end -}}
{{ end -}}
{{ end -}}

{{ range .Versions }}
## {{ .Tag.Name }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ range .CommitGroups -}}
{{ range .Commits -}}
- {{ .Header }}
{{ end }}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits -}}
### Merges
{{ range .MergeCommits -}}
- {{ .Header }}
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
