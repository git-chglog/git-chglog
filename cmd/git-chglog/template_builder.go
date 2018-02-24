package main

import "fmt"

// TemplateBuilder ...
type TemplateBuilder interface {
	Builder
}

type templateBuilderImpl struct{}

// NewTemplateBuilder ...
func NewTemplateBuilder() TemplateBuilder {
	return &templateBuilderImpl{}
}

// Build ...
func (t *templateBuilderImpl) Build(ans *Answer) (string, error) {
	// versions
	tpl := "{{range .Versions}}\n"

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
	tpl += "\n{{end}}"

	return tpl, nil
}

func (*templateBuilderImpl) versionHeader(style, template string) string {
	var (
		tpl     string
		tagName string
		date    = "{{datetime \"2006-01-02\" .Tag.Date}}"
	)

	// parts
	switch style {
	case styleGitHub:
		tpl = "<a name=\"{{.Tag.Name}}\"></a>\n"
		tagName = "{{if .Tag.Previous}}[{{.Tag.Name}}]({{$.Info.RepositoryURL}}/compare/{{.Tag.Previous.Name}}...{{.Tag.Name}}){{else}}{{.Tag.Name}}{{end}}"
	default:
		tagName = "{{.Tag.Name}}"
	}

	// format
	switch template {
	case tplStandard:
		tpl = fmt.Sprintf("%s## %s (%s)\n",
			tpl,
			tagName,
			date,
		)
	case tplCool:
		tpl = fmt.Sprintf("%s## %s\n\n> %s\n",
			tpl,
			tagName,
			date,
		)
	}

	return tpl
}

func (*templateBuilderImpl) commits(template, format string) string {
	var (
		header string
		body   string
	)

	switch format {
	case fmtSubject.Display:
		body = `{{range .Commits}}
* {{.Header}}{{end}}
`

	default:
		if format == fmtTypeScopeSubject.Display {
			header = "{{if ne .Scope \"\"}}**{{.Scope}}:** {{end}}{{.Subject}}"
		} else {
			header = "{{.Subject}}"
		}

		body = fmt.Sprintf(`### {{.Title}}
{{range .Commits}}
* %s{{end}}
`, header)
	}

	return fmt.Sprintf(`{{range .CommitGroups}}
%s{{end}}`, body)
}

func (*templateBuilderImpl) reverts() string {
	return `{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Revert.Header}}{{end}}
{{end}}`
}

func (t *templateBuilderImpl) merges(style string) string {
	var title string

	switch style {
	case styleGitHub:
		title = "Pull Requests"
	default:
		title = "Merges"
	}

	return fmt.Sprintf(`{{if .MergeCommits}}
### %s
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}`, title)
}

func (*templateBuilderImpl) notes() string {
	return `{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}`
}
