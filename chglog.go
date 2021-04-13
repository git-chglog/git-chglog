// Package chglog implements main logic for the CHANGELOG generate.
package chglog

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/tsuyoshiwada/go-gitcmd"
)

// Options is an option used to process commits
type Options struct {
	Processor                   Processor
	NextTag                     string              // Treat unreleased commits as specified tags (EXPERIMENTAL)
	TagFilterPattern            string              // Filter tag by regexp
	Sort                        string              // Specify how to sort tags; currently supports "date" (default) or by "semver".
	NoCaseSensitive             bool                // Filter commits in a case insensitive way
	CommitFilters               map[string][]string // Filter by using `Commit` properties and values. Filtering is not done by specifying an empty value
	CommitSortBy                string              // Property name to use for sorting `Commit` (e.g. `Scope`)
	CommitGroupBy               string              // Property name of `Commit` to be grouped into `CommitGroup` (e.g. `Type`)
	CommitGroupSortBy           string              // Property name to use for sorting `CommitGroup` (e.g. `Title`)
	CommitGroupTitleOrder       []string            // Predefined sorted list of titles to use for sorting `CommitGroup`. Only if `CommitGroupSortBy` is `Custom`
	CommitGroupTitleMaps        map[string]string   // Map for `CommitGroup` title conversion
	HeaderPattern               string              // A regular expression to use for parsing the commit header
	HeaderPatternMaps           []string            // A rule for mapping the result of `HeaderPattern` to the property of `Commit`
	IssuePrefix                 []string            // Prefix used for issues (e.g. `#`, `gh-`)
	RefActions                  []string            // Word list of `Ref.Action`
	MergePattern                string              // A regular expression to use for parsing the merge commit
	MergePatternMaps            []string            // Similar to `HeaderPatternMaps`
	RevertPattern               string              // A regular expression to use for parsing the revert commit
	RevertPatternMaps           []string            // Similar to `HeaderPatternMaps`
	NoteKeywords                []string            // Keyword list to find `Note`. A semicolon is a separator, like `<keyword>:` (e.g. `BREAKING CHANGE`)
	JiraUsername                string
	JiraToken                   string
	JiraURL                     string
	JiraTypeMaps                map[string]string
	JiraIssueDescriptionPattern string
	Paths                       []string // Path filter
}

// Info is metadata related to CHANGELOG
type Info struct {
	Title         string // Title of CHANGELOG
	RepositoryURL string // URL of git repository
}

// RenderData is the data passed to the template
type RenderData struct {
	Info       *Info
	Unreleased *Unreleased
	Versions   []*Version
}

// Config for generating CHANGELOG
type Config struct {
	Bin        string // Git execution command
	WorkingDir string // Working directory
	Template   string // Path for template file. If a relative path is specified, it depends on the value of `WorkingDir`.
	Info       *Info
	Options    *Options
}

func normalizeConfig(config *Config) {
	opts := config.Options

	if opts.HeaderPattern == "" {
		opts.HeaderPattern = "^(.*)$"
		opts.HeaderPatternMaps = []string{
			"Subject",
		}
	}

	if opts.MergePattern == "" {
		opts.MergePattern = "^Merge branch '(\\w+)'$"
		opts.MergePatternMaps = []string{
			"Source",
		}
	}

	if opts.RevertPattern == "" {
		opts.RevertPattern = "^Revert \"([\\s\\S]*)\"$"
		opts.RevertPatternMaps = []string{
			"Header",
		}
	}

	config.Options = opts
}

// Generator of CHANGELOG
type Generator struct {
	client          gitcmd.Client
	config          *Config
	tagReader       *tagReader
	tagSelector     *tagSelector
	commitParser    *commitParser
	commitExtractor *commitExtractor
}

// NewGenerator receives `Config` and create an new `Generator`
func NewGenerator(logger *Logger, config *Config) *Generator {
	client := gitcmd.New(&gitcmd.Config{
		Bin: config.Bin,
	})

	jiraClient := NewJiraClient(config)

	if config.Options.Processor != nil {
		config.Options.Processor.Bootstrap(config)
	}

	normalizeConfig(config)

	return &Generator{
		client:          client,
		config:          config,
		tagReader:       newTagReader(client, config.Options.TagFilterPattern, config.Options.Sort),
		tagSelector:     newTagSelector(),
		commitParser:    newCommitParser(logger, client, jiraClient, config),
		commitExtractor: newCommitExtractor(config.Options),
	}
}

// Generate gets the commit based on the specified tag `query` and writes the result to `io.Writer`
//
// tag `query` can be specified with the following rule
//  <old>..<new> - Commit contained in `<new>` tags from `<old>` (e.g. `1.0.0..2.0.0`)
//  <tagname>..  - Commit from the `<tagname>` to the latest tag (e.g. `1.0.0..`)
//  ..<tagname>  - Commit from the oldest tag to `<tagname>` (e.g. `..1.0.0`)
//  <tagname>    - Commit contained in `<tagname>` (e.g. `1.0.0`)
func (gen *Generator) Generate(w io.Writer, query string) error {
	back, err := gen.workdir()
	if err != nil {
		return err
	}
	defer func() {
		if err = back(); err != nil {
			log.Fatal(err)
		}
	}()

	tags, first, err := gen.getTags(query)
	if err != nil {
		return err
	}

	unreleased, err := gen.readUnreleased(tags)
	if err != nil {
		return err
	}

	versions, err := gen.readVersions(tags, first)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		return fmt.Errorf("commits corresponding to \"%s\" was not found", query)
	}

	return gen.render(w, unreleased, versions)
}

func (gen *Generator) readVersions(tags []*Tag, first string) ([]*Version, error) {
	next := gen.config.Options.NextTag
	versions := []*Version{}

	for i, tag := range tags {
		var (
			isNext = next == tag.Name
			rev    string
		)

		if isNext {
			if tag.Previous != nil {
				rev = tag.Previous.Name + "..HEAD"
			} else {
				rev = "HEAD"
			}
		} else {
			if i+1 < len(tags) {
				rev = tags[i+1].Name + ".." + tag.Name
			} else {
				if first != "" {
					rev = first + ".." + tag.Name
				} else {
					rev = tag.Name
				}
			}
		}

		commits, err := gen.commitParser.Parse(rev)
		if err != nil {
			return nil, err
		}

		commitGroups, mergeCommits, revertCommits, noteGroups := gen.commitExtractor.Extract(commits)

		versions = append(versions, &Version{
			Tag:           tag,
			CommitGroups:  commitGroups,
			Commits:       commits,
			MergeCommits:  mergeCommits,
			RevertCommits: revertCommits,
			NoteGroups:    noteGroups,
		})

		// Instead of `getTags()`, assign the date to the tag
		if isNext && len(commits) != 0 {
			tag.Date = commits[0].Author.Date
		}
	}

	return versions, nil
}

func (gen *Generator) readUnreleased(tags []*Tag) (*Unreleased, error) {
	if gen.config.Options.NextTag != "" {
		return &Unreleased{}, nil
	}

	rev := "HEAD"

	if len(tags) > 0 {
		rev = tags[0].Name + "..HEAD"
	}

	commits, err := gen.commitParser.Parse(rev)
	if err != nil {
		return nil, err
	}

	commitGroups, mergeCommits, revertCommits, noteGroups := gen.commitExtractor.Extract(commits)

	unreleased := &Unreleased{
		CommitGroups:  commitGroups,
		Commits:       commits,
		MergeCommits:  mergeCommits,
		RevertCommits: revertCommits,
		NoteGroups:    noteGroups,
	}

	return unreleased, nil
}

func (gen *Generator) getTags(query string) ([]*Tag, string, error) {
	tags, err := gen.tagReader.ReadAll()
	if err != nil {
		return nil, "", err
	}

	next := gen.config.Options.NextTag
	if next != "" {
		for _, tag := range tags {
			if next == tag.Name {
				return nil, "", fmt.Errorf("\"%s\" tag already exists", next)
			}
		}

		var previous *RelateTag
		if len(tags) > 0 {
			previous = &RelateTag{
				Name:    tags[0].Name,
				Subject: tags[0].Subject,
				Date:    tags[0].Date,
			}
		}

		// Assign the date with `readVersions()`
		tags = append([]*Tag{
			{
				Name:     next,
				Subject:  next,
				Previous: previous,
			},
		}, tags...)
	}

	if len(tags) == 0 {
		return nil, "", errors.New("git-tag does not exist")
	}

	first := ""
	if query != "" {
		tags, first, err = gen.tagSelector.Select(tags, query)
		if err != nil {
			return nil, "", err
		}
	}

	return tags, first, nil
}

func (gen *Generator) workdir() (func() error, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(gen.config.WorkingDir)
	if err != nil {
		return nil, err
	}

	return func() error {
		return os.Chdir(cwd)
	}, nil
}

func (gen *Generator) render(w io.Writer, unreleased *Unreleased, versions []*Version) error {
	if _, err := os.Stat(gen.config.Template); err != nil {
		return err
	}

	fmap := template.FuncMap{
		// format the input time according to layout
		"datetime": func(layout string, input time.Time) string {
			return input.Format(layout)
		},
		// upper case the first character of a string
		"upperFirst": func(s string) string {
			if len(s) > 0 {
				return strings.ToUpper(string(s[0])) + s[1:]
			}
			return ""
		},
		// indent all lines of s n spaces
		"indent": func(s string, n int) string {
			if len(s) == 0 {
				return ""
			}
			pad := strings.Repeat(" ", n)
			return pad + strings.ReplaceAll(s, "\n", "\n"+pad)
		},
		// While Sprig provides these functions, they change the standard input
		// order which leads to a regression. For an example see:
		// https://github.com/Masterminds/sprig/blob/master/functions.go#L149
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
		"replace":   strings.Replace,
	}

	fname := filepath.Base(gen.config.Template)

	t := template.Must(template.New(fname).Funcs(sprig.TxtFuncMap()).Funcs(fmap).ParseFiles(gen.config.Template))

	return t.Execute(w, &RenderData{
		Info:       gen.config.Info,
		Unreleased: unreleased,
		Versions:   versions,
	})
}
