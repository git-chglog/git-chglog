package main

import (
	"fmt"
	"strings"
)

// Previewable ...
type Previewable interface {
	Display() string
	Preview() string
}

// Defaults
var (
	defaultConfigDir        = ".chglog"
	defaultConfigFilename   = "config.yml"
	defaultTemplateFilename = "CHANGELOG.tpl.md"
)

// Styles
var (
	styleGitHub    = "github"
	styleGitLab    = "gitlab"
	styleBitbucket = "bitbucket"
	styleNone      = "none"
	styles         = []string{
		styleGitHub,
		styleGitLab,
		styleBitbucket,
		styleNone,
	}
)

// CommitMessageFormat ...
type CommitMessageFormat struct {
	display     string
	preview     string
	pattern     string
	patternMaps []string
}

// Display ...
func (f *CommitMessageFormat) Display() string {
	return f.display
}

// Preview ...
func (f *CommitMessageFormat) Preview() string {
	return f.preview
}

// PatternMapString ...
func (f *CommitMessageFormat) PatternMapString() string {
	s := " []"
	l := len(f.patternMaps)
	if l == 0 {
		return s
	}

	arr := make([]string, l)
	for i, p := range f.patternMaps {
		arr[i] = fmt.Sprintf(
			"%s- %s",
			strings.Repeat(" ", 6),
			p,
		)
	}

	return fmt.Sprintf("\n%s", strings.Join(arr, "\n"))
}

// Formats
var (
	fmtTypeScopeSubject = &CommitMessageFormat{
		display:     "<type>(<scope>): <subject>",
		preview:     "feat(core): Add new feature",
		pattern:     `^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$`,
		patternMaps: []string{"Type", "Scope", "Subject"},
	}
	fmtTypeSubject = &CommitMessageFormat{
		display:     "<type>: <subject>",
		preview:     "feat: Add new feature",
		pattern:     `^(\\w*)\\:\\s(.*)$`,
		patternMaps: []string{"Type", "Subject"},
	}
	fmtGitBasic = &CommitMessageFormat{
		display:     "<<type> subject>",
		preview:     "Add new feature",
		pattern:     `^((\\w+)\\s.*)$`,
		patternMaps: []string{"Subject", "Type"},
	}
	fmtSubject = &CommitMessageFormat{
		display:     "<subject>",
		preview:     "Add new feature (Not detect `type` field)",
		pattern:     `^(.*)$`,
		patternMaps: []string{"Subject"},
	}
	formats = []Previewable{
		fmtTypeScopeSubject,
		fmtTypeSubject,
		fmtGitBasic,
		fmtSubject,
	}
)

// TemplateStyleFormat ...
type TemplateStyleFormat struct {
	preview string
	display string
}

// Display ...
func (t *TemplateStyleFormat) Display() string {
	return t.display
}

// Preview ...
func (t *TemplateStyleFormat) Preview() string {
	return t.preview
}

// Templates
var (
	tplKeepAChangelog = &TemplateStyleFormat{
		display: "keep-a-changelog",
		preview: "https://github.com/git-chglog/example-type-scope-subject/blob/master/CHANGELOG.kac.md",
	}
	tplStandard = &TemplateStyleFormat{
		display: "standard",
		preview: "https://github.com/git-chglog/example-type-scope-subject/blob/master/CHANGELOG.standard.md",
	}
	tplCool = &TemplateStyleFormat{
		display: "cool",
		preview: "https://github.com/git-chglog/example-type-scope-subject/blob/master/CHANGELOG.cool.md",
	}
	templates = []Previewable{
		tplKeepAChangelog,
		tplStandard,
		tplCool,
	}
)
