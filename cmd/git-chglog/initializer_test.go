package main

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializer(t *testing.T) {
	assert := assert.New(t)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	mockFs := &mockFileSystem{
		ReturnMkdirP: func(path string) error {
			if path == filepath.FromSlash("/test/config") {
				return nil
			}
			return errors.New("")
		},
		ReturnWriteFile: func(path string, content []byte) error {
			if path == filepath.FromSlash("/test/config/config.yml") || path == filepath.FromSlash("/test/config/CHANGELOG.tpl.md") {
				return nil
			}
			return errors.New("")
		},
	}

	questioner := &mockQuestionerImpl{
		ReturnAsk: func() (*Answer, error) {
			return &Answer{
				ConfigDir: "config",
			}, nil
		},
	}

	configBuilder := &mockConfigBuilderImpl{
		ReturnBuild: func(ans *Answer) (string, error) {
			if ans.ConfigDir == "config" {
				return "config", nil
			}
			return "", errors.New("")
		},
	}

	tplBuilder := &mockTemplateBuilderImpl{
		ReturnBuild: func(ans *Answer) (string, error) {
			if ans.ConfigDir == "config" {
				return "template", nil
			}
			return "", errors.New("")
		},
	}

	init := NewInitializer(
		&InitContext{
			WorkingDir: "/test",
			Stdout:     stdout,
			Stderr:     stderr,
		},
		mockFs,
		questioner,
		configBuilder,
		func(t string) TemplateBuilder {
			return tplBuilder
		},
	)

	assert.Equal(ExitCodeOK, init.Run())
	assert.Equal("", stderr.String())
	assert.Contains(stdout.String(), "Configuration file and template")
}
