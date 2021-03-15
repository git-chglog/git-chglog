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

type typeSample struct {
	typeName string
	title    string
}

// CommitMessageFormat ...
type CommitMessageFormat struct {
	display     string
	preview     string
	pattern     string
	patternMaps []string
	typeSamples []typeSample
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
	if len(f.patternMaps) == 0 {
		return " []"
	}

	arr := make([]string, len(f.patternMaps))
	for i, p := range f.patternMaps {
		arr[i] = fmt.Sprintf(
			"%s- %s",
			strings.Repeat(" ", 6),
			p,
		)
	}

	return fmt.Sprintf("\n%s", strings.Join(arr, "\n"))
}

// FilterTypesString ...
func (f *CommitMessageFormat) FilterTypesString() string {
	if len(f.typeSamples) == 0 {
		return " []"
	}

	arr := make([]string, len(f.typeSamples))
	for i, t := range f.typeSamples {
		arr[i] = fmt.Sprintf(
			"%s#%s- %s",
			strings.Repeat(" ", 4), strings.Repeat(" ", 5),
			t.typeName)
	}
	return fmt.Sprintf("\n%s", strings.Join(arr, "\n"))
}

// TitleMapsString ...
func (f *CommitMessageFormat) TitleMapsString() string {
	if len(f.typeSamples) == 0 {
		return " []"
	}

	arr := make([]string, len(f.typeSamples))
	for i, t := range f.typeSamples {
		arr[i] = fmt.Sprintf(
			"%s#%s%s: %s",
			strings.Repeat(" ", 4), strings.Repeat(" ", 3),
			t.typeName, t.title)
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
		typeSamples: []typeSample{
			{"feat", "Features"}, {"fix", "Bug Fixes"},
			{"perf", "Performance Improvements"}, {"refactor", "Code Refactoring"}},
	}
	fmtTypeSubject = &CommitMessageFormat{
		display:     "<type>: <subject>",
		preview:     "feat: Add new feature",
		pattern:     `^(\\w*)\\:\\s(.*)$`,
		patternMaps: []string{"Type", "Subject"},
		typeSamples: []typeSample{
			{"feat", "Features"}, {"fix", "Bug Fixes"},
			{"perf", "Performance Improvements"}, {"refactor", "Code Refactoring"}},
	}
	fmtGitBasic = &CommitMessageFormat{
		display:     "<<type> subject>",
		preview:     "Add new feature",
		pattern:     `^((\\w+)\\s.*)$`,
		patternMaps: []string{"Subject", "Type"},
		typeSamples: []typeSample{
			{"feat", "Features"}, {"fix", "Bug Fixes"},
			{"perf", "Performance Improvements"}, {"refactor", "Code Refactoring"}},
	}
	fmtSubject = &CommitMessageFormat{
		display:     "<subject>",
		preview:     "Add new feature (Not detect `type` field)",
		pattern:     `^(.*)$`,
		patternMaps: []string{"Subject"},
		typeSamples: []typeSample{},
	}
	fmtCommitEmoji = &CommitMessageFormat{
		display:     ":<type>: <subject>",
		preview:     ":sparkles: Add new feature (Commit message with emoji format)",
		pattern:     `^:(\\w*)\\:\\s(.*)$`,
		patternMaps: []string{"Type", "Subject"},
		typeSamples: []typeSample{
			{"sparkles", "Features"}, {"bug", "Bug Fixes"},
			{"zap", "Performance Improvements"}, {"recycle", "Code Refactoring"}},
	}
	formats = []Previewable{
		fmtTypeScopeSubject,
		fmtTypeSubject,
		fmtGitBasic,
		fmtSubject,
		fmtCommitEmoji,
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
