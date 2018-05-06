package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
	survey "gopkg.in/AlecAivazis/survey.v1"
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

// Questioner ...
type Questioner interface {
	Ask() (*Answer, error)
}

type questionerImpl struct {
	client gitcmd.Client
	fs     FileSystem
}

// NewQuestioner ...
func NewQuestioner(client gitcmd.Client, fs FileSystem) Questioner {
	return &questionerImpl{
		client: client,
		fs:     fs,
	}
}

// Ask ...
func (q *questionerImpl) Ask() (*Answer, error) {
	ans, err := q.ask()
	if err != nil {
		return nil, err
	}

	config := filepath.Join(ans.ConfigDir, defaultConfigFilename)
	tpl := filepath.Join(ans.ConfigDir, defaultTemplateFilename)
	c := q.fs.Exists(config)
	t := q.fs.Exists(tpl)
	msg := ""

	if c && t {
		msg = fmt.Sprintf("\"%s\" and \"%s\" already exists. Do you want to overwrite?", config, tpl)
	} else if c {
		msg = fmt.Sprintf("\"%s\" already exists. Do you want to overwrite?", config)
	} else if t {
		msg = fmt.Sprintf("\"%s\" already exists. Do you want to overwrite?", tpl)
	}

	if msg != "" {
		overwrite := false
		err = survey.AskOne(&survey.Confirm{
			Message: msg,
			Default: true,
		}, &overwrite, nil)

		if err != nil || !overwrite {
			return nil, errors.New("creation of the file was interrupted")
		}
	}

	return ans, nil
}

func (q *questionerImpl) ask() (*Answer, error) {
	ans := &Answer{}
	fmts := q.getPreviewableList(formats)
	tpls := q.getPreviewableList(templates)

	var previewableTransform = func(ans interface{}) (newAns interface{}) {
		if s, ok := ans.(string); ok {
			newAns = q.parsePreviewableList(s)
		}
		return
	}

	questions := []*survey.Question{
		{
			Name: "repository_url",
			Prompt: &survey.Input{
				Message: "What is the URL of your repository?",
				Default: q.getRepositoryURL(),
			},
		},
		{
			Name: "style",
			Prompt: &survey.Select{
				Message: "What is your favorite style?",
				Options: styles,
				Default: styles[0],
			},
		},
		{
			Name: "commit_message_format",
			Prompt: &survey.Select{
				Message: "Choose the format of your favorite commit message",
				Options: fmts,
				Default: fmts[0],
			},
			Transform: previewableTransform,
		},
		{
			Name: "template",
			Prompt: &survey.Select{
				Message: "What is your favorite template style?",
				Options: tpls,
				Default: tpls[0],
			},
			Transform: previewableTransform,
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
				Default: defaultConfigDir,
			},
		},
	}

	err := survey.Ask(questions, ans)
	if err != nil {
		return nil, err
	}

	return ans, nil
}

func (*questionerImpl) getPreviewableList(list []Previewable) []string {
	arr := make([]string, len(list))
	max := 0

	for _, p := range list {
		l := len(p.Display())
		if max < l {
			max = l
		}
	}

	for i, p := range list {
		arr[i] = fmt.Sprintf(
			"%s -- %s",
			p.Display()+strings.Repeat(" ", max-len(p.Display())),
			p.Preview(),
		)
	}

	return arr
}

func (*questionerImpl) parsePreviewableList(input string) string {
	return strings.TrimSpace(strings.Split(input, "--")[0])
}

func (q *questionerImpl) getRepositoryURL() string {
	if q.client.CanExec() != nil || q.client.InsideWorkTree() != nil {
		return ""
	}

	rawurl, err := q.client.Exec("config", "--get", "remote.origin.url")
	if err != nil {
		return ""
	}

	return remoteOriginURLToHTTP(rawurl)
}
