package main

import (
	"fmt"
	"strings"
)

// Defaults
var (
	defaultConfigDir        = ".chglog"
	defaultConfigFilename   = "config.yml"
	defaultTemplateFilename = "CHANGELOG.tpl.md"
)

// Styles
var (
	styleGitHub = "github"
	styleGitLab = "gitlab"
	styleNone   = "none"
	styles      = []string{
		styleGitHub,
		styleGitLab,
		styleNone,
	}
)

// CommitMessageFormat ...
type CommitMessageFormat struct {
	Preview     string
	Display     string
	Pattern     string
	PatternMaps []string
}

// PatternMapString ...
func (f *CommitMessageFormat) PatternMapString() string {
	s := " []"
	l := len(f.PatternMaps)
	if l == 0 {
		return s
	}

	arr := make([]string, l)
	for i, p := range f.PatternMaps {
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
		Preview:     "feat(core): Add new feature",
		Display:     "<type>(<scope>): <subject>",
		Pattern:     `^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$`,
		PatternMaps: []string{"Type", "Scope", "Subject"},
	}
	fmtTypeSubject = &CommitMessageFormat{
		Preview:     "feat: Add new feature",
		Display:     "<type>: <subject>",
		Pattern:     `^(\\w*)\\:\\s(.*)$`,
		PatternMaps: []string{"Type", "Subject"},
	}
	fmtGitBasic = &CommitMessageFormat{
		Preview:     "Add new feature",
		Display:     "<<type> subject>",
		Pattern:     `^((\\w+)\\s.*)$`,
		PatternMaps: []string{"Subject", "Type"},
	}
	fmtSubject = &CommitMessageFormat{
		Preview:     "Add new feature (Not detect `type` field)",
		Display:     "<subject>",
		Pattern:     `^(.*)$`,
		PatternMaps: []string{"Subject"},
	}
	formats = []*CommitMessageFormat{
		fmtTypeScopeSubject,
		fmtTypeSubject,
		fmtGitBasic,
		fmtSubject,
	}
)

// Templates
var (
	tplKeepAChangelog = "keep-a-changelog"
	tplStandard       = "standard"
	tplCool           = "cool"
	templates         = []string{
		tplKeepAChangelog,
		tplStandard,
		tplCool,
	}
)
