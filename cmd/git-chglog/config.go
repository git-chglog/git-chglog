package main

import (
	"path/filepath"
	"strings"

	chglog "github.com/git-chglog/git-chglog"
	"github.com/imdario/mergo"
)

// Info ...
type Info struct {
	Title         string `yaml:"title"`
	RepositoryURL string `yaml:"repository_url"`
}

// CommitOptions ...
type CommitOptions struct {
	Filters map[string][]string `yaml:"filters"`
	SortBy  string              `yaml:"sort_by"`
}

// CommitGroupOptions ...
type CommitGroupOptions struct {
	GroupBy   string            `yaml:"group_by"`
	SortBy    string            `yaml:"sort_by"`
	TitleMaps map[string]string `yaml:"title_maps"`
}

// PatternOptions ...
type PatternOptions struct {
	Pattern     string   `yaml:"pattern"`
	PatternMaps []string `yaml:"pattern_maps"`
}

// IssueOptions ...
type IssueOptions struct {
	Prefix []string `yaml:"prefix"`
}

// RefOptions ...
type RefOptions struct {
	Actions []string `yaml:"actions"`
}

// NoteOptions ...
type NoteOptions struct {
	Keywords []string `yaml:"keywords"`
}

// Options ...
type Options struct {
	Commits      CommitOptions      `yaml:"commits"`
	CommitGroups CommitGroupOptions `yaml:"commit_groups"`
	Header       PatternOptions     `yaml:"header"`
	Issues       IssueOptions       `yaml:"issues"`
	Refs         RefOptions         `yaml:"refs"`
	Merges       PatternOptions     `yaml:"merges"`
	Reverts      PatternOptions     `yaml:"reverts"`
	Notes        NoteOptions        `yaml:"notes"`
}

// Config ...
type Config struct {
	Bin      string  `yaml:"bin"`
	Template string  `yaml:"template"`
	Style    string  `yaml:"style"`
	Info     Info    `yaml:"info"`
	Options  Options `yaml:"options"`
}

// Normalize ...
func (config *Config) Normalize(ctx *CLIContext) error {
	err := mergo.Merge(config, &Config{
		Bin:      "git",
		Template: "CHANGELOG.tpl.md",
		Info: Info{
			Title: "CHANGELOG",
		},
		Options: Options{
			Commits: CommitOptions{
				SortBy: "Scope",
			},
			CommitGroups: CommitGroupOptions{
				GroupBy: "Type",
				SortBy:  "Title",
			},
		},
	})

	if err != nil {
		return err
	}

	config.Info.RepositoryURL = strings.TrimRight(config.Info.RepositoryURL, "/")

	if !filepath.IsAbs(config.Template) {
		config.Template = filepath.Join(filepath.Dir(ctx.ConfigPath), config.Template)
	}

	config.normalizeStyle()

	return nil
}

// Normalize style
func (config *Config) normalizeStyle() {
	switch config.Style {
	case "github":
		config.normalizeStyleOfGitHub()
	case "gitlab":
		config.normalizeStyleOfGitLab()
	case "bitbucket":
		config.normalizeStyleOfBitbucket()
	}
}

// For GitHub
func (config *Config) normalizeStyleOfGitHub() {
	opts := config.Options

	if len(opts.Issues.Prefix) == 0 {
		opts.Issues.Prefix = []string{
			"#",
			"gh-",
		}
	}

	if len(opts.Refs.Actions) == 0 {
		opts.Refs.Actions = []string{
			"close",
			"closes",
			"closed",
			"fix",
			"fixes",
			"fixed",
			"resolve",
			"resolves",
			"resolved",
		}
	}

	if opts.Merges.Pattern == "" && len(opts.Merges.PatternMaps) == 0 {
		opts.Merges.Pattern = "^Merge pull request #(\\d+) from (.*)$"
		opts.Merges.PatternMaps = []string{
			"Ref",
			"Source",
		}
	}

	config.Options = opts
}

// For GitLab
func (config *Config) normalizeStyleOfGitLab() {
	opts := config.Options

	if len(opts.Issues.Prefix) == 0 {
		opts.Issues.Prefix = []string{
			"#",
		}
	}

	if len(opts.Refs.Actions) == 0 {
		opts.Refs.Actions = []string{
			"close",
			"closes",
			"closed",
			"closing",
			"fix",
			"fixes",
			"fixed",
			"fixing",
			"resolve",
			"resolves",
			"resolved",
			"resolving",
		}
	}

	if opts.Merges.Pattern == "" && len(opts.Merges.PatternMaps) == 0 {
		opts.Merges.Pattern = "^Merge branch '.*' into '(.*)'$"
		opts.Merges.PatternMaps = []string{
			"Source",
		}
	}

	config.Options = opts
}

// For Bitbucket
func (config *Config) normalizeStyleOfBitbucket() {
	opts := config.Options

	if len(opts.Issues.Prefix) == 0 {
		opts.Issues.Prefix = []string{
			"#",
		}
	}

	if len(opts.Refs.Actions) == 0 {
		opts.Refs.Actions = []string{
			"close",
			"closes",
			"closed",
			"closing",
			"fix",
			"fixed",
			"fixes",
			"fixing",
			"resolve",
			"resolves",
			"resolved",
			"resolving",
			"eopen",
			"reopens",
			"reopening",
			"hold",
			"holds",
			"holding",
			"wontfix",
			"invalidate",
			"invalidates",
			"invalidated",
			"invalidating",
			"addresses",
			"re",
			"references",
			"ref",
			"refs",
			"see",
		}
	}

	if opts.Merges.Pattern == "" && len(opts.Merges.PatternMaps) == 0 {
		opts.Merges.Pattern = "^Merged in (.*) \\(pull request #(\\d+)\\)$"
		opts.Merges.PatternMaps = []string{
			"Source",
			"Ref",
		}
	}

	config.Options = opts
}

// Convert ...
func (config *Config) Convert(ctx *CLIContext) *chglog.Config {
	info := config.Info
	opts := config.Options

	return &chglog.Config{
		Bin:        config.Bin,
		WorkingDir: ctx.WorkingDir,
		Template:   config.Template,
		Info: &chglog.Info{
			Title:         info.Title,
			RepositoryURL: info.RepositoryURL,
		},
		Options: &chglog.Options{
			NextTag:              ctx.NextTag,
			CommitFilters:        opts.Commits.Filters,
			CommitSortBy:         opts.Commits.SortBy,
			CommitGroupBy:        opts.CommitGroups.GroupBy,
			CommitGroupSortBy:    opts.CommitGroups.SortBy,
			CommitGroupTitleMaps: opts.CommitGroups.TitleMaps,
			HeaderPattern:        opts.Header.Pattern,
			HeaderPatternMaps:    opts.Header.PatternMaps,
			IssuePrefix:          opts.Issues.Prefix,
			RefActions:           opts.Refs.Actions,
			MergePattern:         opts.Merges.Pattern,
			MergePatternMaps:     opts.Merges.PatternMaps,
			RevertPattern:        opts.Reverts.Pattern,
			RevertPatternMaps:    opts.Reverts.PatternMaps,
			NoteKeywords:         opts.Notes.Keywords,
		},
	}
}
