package chglog

import (
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

// Options ...
type Options struct {
	Processor            Processor
	CommitFilters        map[string][]string
	CommitSortBy         string
	CommitGroupBy        string
	CommitGroupSortBy    string
	CommitGroupTitleMaps map[string]string
	HeaderPattern        string
	HeaderPatternMaps    []string
	IssuePrefix          []string
	RefActions           []string
	MergePattern         string
	MergePatternMaps     []string
	RevertPattern        string
	RevertPatternMaps    []string
	NoteKeywords         []string
}

// Info ...
type Info struct {
	Title         string
	RepositoryURL string
}

// RenderData ...
type RenderData struct {
	Info     *Info
	Versions []*Version
}

// Config ...
type Config struct {
	Bin      string
	Path     string
	Template string
	Info     *Info
	Options  *Options
}

// Generator ...
type Generator struct {
	client          gitcmd.Client
	config          *Config
	tagReader       *tagReader
	tagSelector     *tagSelector
	commitParser    *commitParser
	commitExtractor *commitExtractor
}

// NewGenerator ...
func NewGenerator(config *Config) *Generator {
	client := gitcmd.New(&gitcmd.Config{
		Bin: config.Bin,
	})

	if config.Options.Processor != nil {
		config.Options.Processor.Bootstrap(config)
	}

	return &Generator{
		client:          client,
		config:          config,
		tagReader:       newTagReader(client),
		tagSelector:     newTagSelector(),
		commitParser:    newCommitParser(client, config),
		commitExtractor: newCommitExtractor(config.Options),
	}
}

// Generate ...
func (gen *Generator) Generate(w io.Writer, query string) error {
	back, err := gen.workdir()
	if err != nil {
		return err
	}

	versions, err := gen.readVersions(query)
	if err != nil {
		return err
	}

	back()

	return gen.render(w, versions)
}

func (gen *Generator) readVersions(query string) ([]*Version, error) {
	tags, err := gen.tagReader.ReadAll()
	if err != nil {
		return nil, err
	}

	first := ""
	if query != "" {
		tags, first, err = gen.tagSelector.Select(tags, query)
		if err != nil {
			return nil, err
		}
	}

	versions := []*Version{}

	for i, tag := range tags {
		var rev string

		if i+1 < len(tags) {
			rev = tags[i+1].Name + ".." + tag.Name
		} else {
			if first != "" {
				rev = first + ".." + tag.Name
			} else {
				rev = tag.Name
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
			MergeCommits:  mergeCommits,
			RevertCommits: revertCommits,
			NoteGroups:    noteGroups,
		})
	}

	return versions, nil
}

func (gen *Generator) workdir() (func() error, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(gen.config.Path)
	if err != nil {
		return nil, err
	}

	return func() error {
		return os.Chdir(cwd)
	}, nil
}

func (gen *Generator) render(w io.Writer, versions []*Version) error {
	fmap := template.FuncMap{
		"datetime": func(layout string, input time.Time) string {
			return input.Format(layout)
		},
	}

	fname := filepath.Base(gen.config.Template)

	t := template.Must(template.New(fname).Funcs(fmap).ParseFiles(gen.config.Template))

	return t.Execute(w, &RenderData{
		Info:     gen.config.Info,
		Versions: versions,
	})
}
