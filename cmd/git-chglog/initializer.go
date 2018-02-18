package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
	survey "gopkg.in/AlecAivazis/survey.v1"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

var (
	defaultConfigFilename   = "config.yml"
	defaultTemplateFilename = "CHANGELOG.tpl.md"

	styleGitHub     = "github"
	styleNone       = "none"
	changelogStyles = []string{
		styleGitHub,
		styleNone,
	}

	fmtTypeScopeSubject  = "<type>(<scope>): <subject>"
	fmtTypeSubject       = "<type>: <subject>"
	fmtSubject           = "<subject>"
	commitMessageFormats = []string{
		fmtTypeScopeSubject,
		fmtTypeSubject,
		fmtSubject,
	}

	tplStandard    = "standard"
	tplCool        = "cool"
	templateStyles = []string{
		tplStandard,
		tplCool,
	}
)

// Answer ...
type Answer struct {
	RepositoryURL       string `survey:"repository_url"`
	Style               string `survey:"style"`
	CommitMessageFormat string `survey:"commit_message_format"`
	Template            string `survey:"template"`
	IncludeMerges       bool   `survey:"include_merges"`
	IncludeReverts      bool   `survey:"include_reverts"`
	ConfigDir           string `survey:"config_dir"`
}

// Initializer ...
type Initializer struct {
	client gitcmd.Client
}

// NewInitializer ...
func NewInitializer() *Initializer {
	return &Initializer{
		client: gitcmd.New(&gitcmd.Config{
			Bin: "git",
		}),
	}
}

// Run ...
func (init *Initializer) Run() int {
	answer, err := init.ask()
	if err != nil {
		return ExitCodeError
	}

	err = init.generateConfigure(answer)
	if err != nil {
		return ExitCodeError
	}

	success := color.CyanString("✔")
	emoji.Fprintf(os.Stdout, `
:sparkles:%s
  %s %s
  %s %s

`,
		color.GreenString("Configuration file and template generation completed!"),
		success,
		filepath.Join(answer.ConfigDir, defaultConfigFilename),
		success,
		filepath.Join(answer.ConfigDir, defaultTemplateFilename),
	)

	return ExitCodeOK
}

func (init *Initializer) ask() (*Answer, error) {
	answer := &Answer{}
	qs := init.createQuestions()
	err := survey.Ask(qs, answer)
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (init *Initializer) createQuestions() []*survey.Question {
	originURL := init.getRepositoryURL()

	return []*survey.Question{
		{
			Name: "repository_url",
			Prompt: &survey.Input{
				Message: "What is the URL of your repository?",
				Default: originURL,
			},
		},
		{
			Name: "style",
			Prompt: &survey.Select{
				Message: "What is your favorite style?",
				Options: changelogStyles,
				Default: changelogStyles[0],
			},
		},
		{
			Name: "commit_message_format",
			Prompt: &survey.Select{
				Message: "Choose the format of your favorite commit message",
				Options: commitMessageFormats,
				Default: commitMessageFormats[0],
			},
		},
		{
			Name: "template",
			Prompt: &survey.Select{
				Message: "What is your favorite template style?",
				Options: templateStyles,
				Default: templateStyles[0],
			},
		},
		{
			Name: "include_merges",
			Prompt: &survey.Confirm{
				Message: "Do you include Merge Commit in CHANGELOG?",
				Default: true,
			},
		},
		{
			Name: "include_reverts",
			Prompt: &survey.Confirm{
				Message: "Do you include Revert Commit in CHANGELOG?",
				Default: true,
			},
		},
		{
			Name: "config_dir",
			Prompt: &survey.Input{
				Message: "In which directory do you output configuration files and templates?",
				Default: ".chglog",
			},
		},
	}
}

func (init *Initializer) getRepositoryURL() string {
	if init.client.CanExec() != nil || init.client.InsideWorkTree() != nil {
		return ""
	}

	rawurl, err := init.client.Exec("config", "--get", "remote.origin.url")
	if err != nil {
		return ""
	}

	return remoteOriginURLToHTTP(rawurl)
}

func (init *Initializer) generateConfigure(answer *Answer) error {
	var err error

	err = fs.MkdirP(answer.ConfigDir)
	if err != nil {
		return err
	}

	config := init.createConfigYamlContent(answer)
	tpl := init.createTemplate(answer)

	configPath := filepath.Join(answer.ConfigDir, defaultConfigFilename)
	templatePath := filepath.Join(answer.ConfigDir, defaultTemplateFilename)

	err = init.createFileWithConfirm(configPath, config)
	if err != nil {
		return err
	}

	err = init.createFileWithConfirm(templatePath, tpl)
	if err != nil {
		return err
	}

	return nil
}

func (*Initializer) createFileWithConfirm(path, content string) error {
	if _, err := os.Stat(path); err == nil {
		answer := struct {
			OK bool
		}{}

		err := survey.Ask([]*survey.Question{
			{
				Name: "ok",
				Prompt: &survey.Confirm{
					Message: fmt.Sprintf("\"%s\" already exists. Do you want to overwrite?", path),
					Default: true,
				},
			},
		}, &answer)

		if err != nil {
			return err
		}

		if !answer.OK {
			return errors.New("creation of the file was interrupted")
		}
	}

	err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (init *Initializer) createConfigYamlContent(answer *Answer) string {
	var (
		style             = answer.Style
		template          = defaultTemplateFilename
		repositoryURL     = answer.RepositoryURL
		headerPattern     string
		headerPatternMaps string
	)

	switch answer.CommitMessageFormat {
	case fmtTypeScopeSubject:
		headerPattern = `^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$`
		headerPatternMaps = `
      - Type
      - Scope
      - Subject`
	case fmtTypeSubject:
		headerPattern = `^(\\w*)\\:\\s(.*)$`
		headerPatternMaps = `
      - Type
      - Subject`
	case fmtSubject:
		headerPattern = `^(.*)$`
		headerPatternMaps = `
      - Subject`
	}

	config := fmt.Sprintf(`style: %s
template: %s
info:
  title: CHANGELOG
  repository_url: %s
options:
  commits:
    # filters:
    #   Type:
    #     - feat
    #     - fix
    #     - perf
    #     - refactor
  commit_groups:
    # title_maps:
    #   feat: Features
    #   fix: Bug Fixes
    #   perf: Performance Improvements
    #   refactor: Code Refactoring
  header:
    pattern: "%s"
    pattern_maps:%s
  notes:
    keywords:
      - BREAKING CHANGE`,
		style,
		template,
		repositoryURL,
		headerPattern,
		headerPatternMaps,
	)

	return config
}

func (init *Initializer) createTemplate(answer *Answer) string {
	tpl := "{{range .Versions}}\n"

	// versions
	tpl += init.versionHeader(answer.Style, answer.Template)

	// commits
	tpl += init.commits(answer.Style, answer.Template, answer.CommitMessageFormat)

	// merges
	if answer.IncludeReverts {
		tpl += `{{if .RevertCommits}}
### Reverts
{{range .RevertCommits}}
* {{.Header}}{{end}}
{{end}}`
	}

	// reverts
	if answer.IncludeReverts {
		tpl += fmt.Sprintf(`{{if .MergeCommits}}
### %s
{{range .MergeCommits}}
* {{.Header}}{{end}}
{{end}}`, init.mergeTitle(answer.Style))
	}

	tpl += `{{range .NoteGroups}}
### {{.Title}}
{{range .Notes}}
{{.Body}}
{{end}}
{{end}}
{{end}}`

	return tpl
}

func (*Initializer) versionHeader(style, template string) string {
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

func (*Initializer) commits(style, template, format string) string {
	var (
		header string
		body   string
	)

	switch format {
	case fmtTypeScopeSubject, fmtTypeSubject:
		if format == fmtTypeScopeSubject {
			header = "{{if ne .Scope \"\"}}**{{.Scope}}:** {{end}}{{.Subject}}"
		} else {
			header = "{{.Subject}}"
		}

		body = fmt.Sprintf(`### {{.Title}}
{{range .Commits}}
* %s{{end}}
`, header)

	case fmtSubject:
		body = `{{range .Commits}}
* {{.Header}}{{end}}
`
	}

	return fmt.Sprintf(`{{range .CommitGroups}}
%s{{end}}`, body)
}

func (*Initializer) mergeTitle(style string) string {
	switch style {
	case styleGitHub:
		return "Pull Requests"
	default:
		return "Merges"
	}
}
