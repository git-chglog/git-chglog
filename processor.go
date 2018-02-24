package chglog

import (
	"regexp"
	"strings"
)

// Processor hooks the internal processing of `Generator`, it is possible to adjust the contents
type Processor interface {
	Bootstrap(*Config)
	ProcessCommit(*Commit) *Commit
}

// GitHubProcessor is optimized for CHANGELOG used in GitHub
//
// The following processing is performed
//  - Mentions automatic link (@tsuyoshiwada -> [@tsuyoshiwada](https://github.com/tsuyoshiwada))
//  - Automatic link to references (#123 -> [#123](https://github.com/owner/repo/issues/123))
type GitHubProcessor struct {
	Host      string // Host name used for link destination. Note: You must include the protocol (e.g. "https://github.com")
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

	if commit.Revert != nil {
		commit.Revert.Header = p.addLinks(commit.Revert.Header)
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
