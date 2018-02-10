package chglog

import (
	"regexp"
	"strings"
)

// Processor ...
type Processor interface {
	Bootstrap(*Config)
	ProcessCommit(*Commit) *Commit
}

// PlainProcessor ...
type PlainProcessor struct{}

// Bootstrap ...
func (*PlainProcessor) Bootstrap(config *Config) {}

// ProcessCommit ...
func (*PlainProcessor) ProcessCommit(commit *Commit) *Commit {
	return commit
}

// GitHubProcessor ...
type GitHubProcessor struct {
	Host      string
	config    *Config
	reMention *regexp.Regexp
	reIssue   *regexp.Regexp
}

// Bootstrap ...
func (p *GitHubProcessor) Bootstrap(config *Config) {
	p.config = config

	if p.Host == "" {
		p.Host = "https://github.com"
	} else {
		p.Host = strings.TrimRight(p.Host, "/")
	}

	p.reMention = regexp.MustCompile("@(\\w+)")
	p.reIssue = regexp.MustCompile("(?i)(#|gh-)(\\d+)")
}

// ProcessCommit ...
func (p *GitHubProcessor) ProcessCommit(commit *Commit) *Commit {
	commit.Header = p.addLinks(commit.Header)
	commit.Subject = p.addLinks(commit.Subject)
	commit.Body = p.addLinks(commit.Body)
	for _, note := range commit.Notes {
		note.Body = p.addLinks(note.Body)
	}
	return commit
}

func (p *GitHubProcessor) addLinks(input string) string {
	repoURL := strings.TrimRight(p.config.Info.RepositoryURL, "/")

	// mentions
	input = p.reMention.ReplaceAllString(input, "[@$1]("+p.Host+"/$1)")

	// issues
	input = p.reIssue.ReplaceAllString(input, "[$1$2]("+repoURL+"/issues/$2)")

	return input
}
