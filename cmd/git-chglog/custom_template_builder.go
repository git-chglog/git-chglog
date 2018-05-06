package main

import "fmt"

type customTemplateBuilderImpl struct{}

// NewCustomTemplateBuilder ...
func NewCustomTemplateBuilder() TemplateBuilder {
	return &customTemplateBuilderImpl{}
}

// Build ...
func (t *customTemplateBuilderImpl) Build(ans *Answer) (string, error) {
	// versions
	tpl := "{{ range .Versions }}\n"

	// version header
	tpl += t.versionHeader(ans.Style, ans.Template)

	// commits
	tpl += t.commits(ans.Template, ans.CommitMessageFormat)

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
	tpl += "{{ end -}}"

	return tpl, nil
}

func (*customTemplateBuilderImpl) versionHeader(style, template string) string {
	var (
		tpl     string
		tagName = "{{ .Tag.Name }}"
		date    = "{{ datetime \"2006-01-02\" .Tag.Date }}"
	)

	// parts
	switch style {
	case styleGitHub, styleGitLab:
		tpl = templateTagNameAnchor
		tagName = "{{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }}"
	case styleBitbucket:
		tpl = templateTagNameAnchor
		tagName = "{{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Name }}..{{ .Tag.Previous.Name }}){{ else }}{{ .Tag.Name }}{{ end }}"
	}

	// format
	switch template {
	case tplStandard.display:
		tpl = fmt.Sprintf("%s## %s (%s)\n\n",
			tpl,
			tagName,
			date,
		)
	case tplCool.display:
		tpl = fmt.Sprintf("%s## %s\n\n> %s\n\n",
			tpl,
			tagName,
			date,
		)
	}

	return tpl
}

func (*customTemplateBuilderImpl) commits(template, format string) string {
	var (
		header string
		body   string
	)

	switch format {
	case fmtSubject.display:
		body = `{{ range .Commits -}}
* {{ .Header }}
{{ end }}`

	default:
		if format == fmtTypeScopeSubject.display {
			header = "{{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}"
		} else {
			header = "{{ .Subject }}"
		}

		body = fmt.Sprintf(`### {{ .Title }}

{{ range .Commits -}}
* %s
{{ end }}`, header)
	}

	return fmt.Sprintf(`{{ range .CommitGroups -}}
%s
{{ end -}}
`, body)
}

func (*customTemplateBuilderImpl) reverts() string {
	return `
{{- if .RevertCommits -}}
### Reverts

{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}
`
}

func (t *customTemplateBuilderImpl) merges(style string) string {
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
* {{ .Header }}
{{ end }}
{{ end -}}
`, title)
}

func (*customTemplateBuilderImpl) notes() string {
	return `
{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}

{{ range .Notes }}
{{ .Body }}
{{ end }}
{{ end -}}
{{ end -}}
`
}
