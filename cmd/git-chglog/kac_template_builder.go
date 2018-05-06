package main

import "fmt"

type kacTemplateBuilderImpl struct{}

// NewKACTemplateBuilder ...
func NewKACTemplateBuilder() TemplateBuilder {
	return &kacTemplateBuilderImpl{}
}

// Build ...
func (t *kacTemplateBuilderImpl) Build(ans *Answer) (string, error) {
	// unreleased
	tpl := t.unreleased(ans.Style, ans.CommitMessageFormat)

	// version start
	tpl += "\n{{ range .Versions }}\n"

	tpl += t.versionHeader(ans.Style)

	// commits
	tpl += t.commits(".CommitGroups", ans.CommitMessageFormat)

	// revert
	if ans.IncludeReverts {
		tpl += t.reverts()
	}

	// merges
	if ans.IncludeMerges {
		tpl += t.merges(ans.Style)
	}

	// notes
	tpl += t.notes()

	// versions end
	tpl += "\n{{ end -}}"

	// footer (links)
	tpl += t.footer(ans.Style)

	return tpl, nil
}

func (t *kacTemplateBuilderImpl) unreleased(style, format string) string {
	var (
		id      = ""
		title   = "Unreleased"
		commits = t.commits(".Unreleased.CommitGroups", format)
	)

	switch style {
	case styleGitHub, styleGitLab, styleBitbucket:
		id = "<a name=\"unreleased\"></a>\n"
		title = fmt.Sprintf("[%s]", title)
	}

	return fmt.Sprintf(`{{ if .Versions -}}
%s## %s

{{ if .Unreleased.CommitGroups -}}
%s{{ end -}}
{{ end -}}
`, id, title, commits)
}

func (t *kacTemplateBuilderImpl) versionHeader(style string) string {
	var (
		id      = ""
		tagName = "{{ .Tag.Name }}"
		date    = "{{ datetime \"2006-01-02\" .Tag.Date }}"
	)

	switch style {
	case styleGitHub, styleGitLab, styleBitbucket:
		id = templateTagNameAnchor
		tagName = "{{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }}"
	}

	return fmt.Sprintf("%s## %s - %s\n", id, tagName, date)
}

func (t *kacTemplateBuilderImpl) commits(commitGroups, format string) string {
	var (
		body string
	)

	switch format {
	case fmtSubject.display:
		body = `{{ range .Commits -}}
- {{ .Header }}
{{ end }}`

	default:
		body = `### {{ .Title }}
{{ range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end }}`
	}

	return fmt.Sprintf(`{{ range %s -}}
%s
{{ end -}}
`, commitGroups, body)
}

func (t *kacTemplateBuilderImpl) reverts() string {
	return `
{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end }}
{{ end -}}
`
}

func (t *kacTemplateBuilderImpl) merges(style string) string {
	var title string

	switch style {
	case styleGitHub, styleBitbucket:
		title = "Pull Requests"
	case styleGitLab:
		title = "Merge Requests"
	default:
		title = "Merges"
	}

	return fmt.Sprintf(`
{{- if .MergeCommits -}}
### %s
{{ range .MergeCommits -}}
- {{ .Header }}
{{ end }}
{{ end -}}
`, title)
}

func (*kacTemplateBuilderImpl) notes() string {
	return `
{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}`
}

func (*kacTemplateBuilderImpl) footer(style string) string {
	switch style {
	case styleGitHub, styleGitLab:
		return `

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ end -}}
{{ end -}}`
	case styleBitbucket:
		return `

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/HEAD..{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Name }}..{{ .Tag.Previous.Name }}
{{ end -}}
{{ end -}}
{{ end -}}`
	default:
		return ""
	}
}
