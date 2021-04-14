package main

import (
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"

	chglog "github.com/git-chglog/git-chglog"
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
	GroupBy    string            `yaml:"group_by"`
	SortBy     string            `yaml:"sort_by"`
	TitleOrder []string          `yaml:"title_order"`
	TitleMaps  map[string]string `yaml:"title_maps"`
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

// JiraClientInfoOptions ...
type JiraClientInfoOptions struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
	URL      string `yaml:"url"`
}

// JiraIssueOptions ...
type JiraIssueOptions struct {
	TypeMaps           map[string]string `yaml:"type_maps"`
	DescriptionPattern string            `yaml:"description_pattern"`
}

// JiraOptions ...
type JiraOptions struct {
	ClintInfo JiraClientInfoOptions `yaml:"info"`
	Issue     JiraIssueOptions      `yaml:"issue"`
}

// Options ...
type Options struct {
	TagFilterPattern string             `yaml:"tag_filter_pattern"`
	Sort             string             `yaml:"sort"`
	Commits          CommitOptions      `yaml:"commits"`
	CommitGroups     CommitGroupOptions `yaml:"commit_groups"`
	Header           PatternOptions     `yaml:"header"`
	Issues           IssueOptions       `yaml:"issues"`
	Refs             RefOptions         `yaml:"refs"`
	Merges           PatternOptions     `yaml:"merges"`
	Reverts          PatternOptions     `yaml:"reverts"`
	Notes            NoteOptions        `yaml:"notes"`
	Jira             JiraOptions        `yaml:"jira"`
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
	config.normalizeTagSortBy()

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

func (config *Config) normalizeTagSortBy() {
	switch {
	case config.Options.Sort == "":
		config.Options.Sort = "date"
	case strings.EqualFold(config.Options.Sort, "date"):
		config.Options.Sort = "date"
	case strings.EqualFold(config.Options.Sort, "semver"):
		config.Options.Sort = "semver"
	default:
		config.Options.Sort = "date"
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

func orValue(str1 string, str2 string) string {
	if str1 != "" {
		return str1
	}
	return str2
}

// Convert ...
func (config *Config) Convert(ctx *CLIContext) *chglog.Config {
	info := config.Info
	opts := config.Options

	if ctx.TagFilterPattern == "" {
		ctx.TagFilterPattern = opts.TagFilterPattern
	}

	return &chglog.Config{
		Bin:        config.Bin,
		WorkingDir: ctx.WorkingDir,
		Template:   orValue(ctx.Template, config.Template),
		Info: &chglog.Info{
			Title:         info.Title,
			RepositoryURL: orValue(ctx.RepositoryURL, info.RepositoryURL),
		},
		Options: &chglog.Options{
			NextTag:                     ctx.NextTag,
			TagFilterPattern:            ctx.TagFilterPattern,
			Sort:                        orValue(ctx.Sort, opts.Sort),
			NoCaseSensitive:             ctx.NoCaseSensitive,
			Paths:                       ctx.Paths,
			CommitFilters:               opts.Commits.Filters,
			CommitSortBy:                opts.Commits.SortBy,
			CommitGroupBy:               opts.CommitGroups.GroupBy,
			CommitGroupSortBy:           opts.CommitGroups.SortBy,
			CommitGroupTitleMaps:        opts.CommitGroups.TitleMaps,
			CommitGroupTitleOrder:       opts.CommitGroups.TitleOrder,
			HeaderPattern:               opts.Header.Pattern,
			HeaderPatternMaps:           opts.Header.PatternMaps,
			IssuePrefix:                 opts.Issues.Prefix,
			RefActions:                  opts.Refs.Actions,
			MergePattern:                opts.Merges.Pattern,
			MergePatternMaps:            opts.Merges.PatternMaps,
			RevertPattern:               opts.Reverts.Pattern,
			RevertPatternMaps:           opts.Reverts.PatternMaps,
			NoteKeywords:                opts.Notes.Keywords,
			JiraUsername:                orValue(ctx.JiraUsername, opts.Jira.ClintInfo.Username),
			JiraToken:                   orValue(ctx.JiraToken, opts.Jira.ClintInfo.Token),
			JiraURL:                     orValue(ctx.JiraURL, opts.Jira.ClintInfo.URL),
			JiraTypeMaps:                opts.Jira.Issue.TypeMaps,
			JiraIssueDescriptionPattern: opts.Jira.Issue.DescriptionPattern,
		},
	}
}
