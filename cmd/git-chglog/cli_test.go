package main

import (
	"bytes"
	"errors"
	"io"
	"testing"

	chglog "github.com/git-chglog/git-chglog"
	"github.com/stretchr/testify/assert"
)

func TestCLIForStdout(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	mockFS := &mockFileSystem{}

	configLoader := &mockConfigLoaderImpl{
		ReturnLoad: func(path string) (*Config, error) {
			if path != "/.chglog/config.yml" {
				return nil, errors.New("")
			}
			return &Config{
				Bin: "/custom/bin/git",
			}, nil
		},
	}

	generator := &mockGeneratorImpl{
		ReturnGenerate: func(w io.Writer, query string, config *chglog.Config) error {
			if config.Bin != "/custom/bin/git" {
				return errors.New("")
			}
			w.Write([]byte("success!!"))
			return nil
		},
	}

	c := NewCLI(
		&Context{
			WorkingDir: "/",
			ConfigPath: "/.chglog/config.yml",
			OutputPath: "",
			Stdout:     stdout,
			Stderr:     stderr,
		},
		mockFS,
		configLoader,
		generator,
	)

	assert.Equal(ExitCodeOK, c.Run())
	assert.Equal("", stderr.String())
	assert.Equal("success!!", stdout.String())
}

func TestCLIForFile(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	mockFS := &mockFileSystem{
		ReturnMkdirP: func(path string) error {
			if path != "/dir/to" {
				return errors.New("")
			}
			return nil
		},
		ReturnCreate: func(name string) (File, error) {
			if name != "/dir/to/CHANGELOG.tpl" {
				return nil, errors.New("")
			}
			return &mockFile{
				ReturnWrite: func(b []byte) (int, error) {
					if string(b) != "success!!" {
						return 0, errors.New("")
					}
					return 0, nil
				},
			}, nil
		},
	}

	configLoader := &mockConfigLoaderImpl{
		ReturnLoad: func(path string) (*Config, error) {
			if path != "/.chglog/config.yml" {
				return nil, errors.New("")
			}
			return &Config{
				Bin: "/custom/bin/git",
			}, nil
		},
	}

	generator := &mockGeneratorImpl{
		ReturnGenerate: func(w io.Writer, query string, config *chglog.Config) error {
			if config.Bin != "/custom/bin/git" {
				return errors.New("")
			}
			w.Write([]byte("success!!"))
			return nil
		},
	}

	c := NewCLI(
		&Context{
			WorkingDir: "/",
			ConfigPath: "/.chglog/config.yml",
			OutputPath: "/dir/to/CHANGELOG.tpl",
			Stdout:     stdout,
			Stderr:     stderr,
		},
		mockFS,
		configLoader,
		generator,
	)

	assert.Equal(ExitCodeOK, c.Run())
	assert.Equal("", stderr.String())
	assert.Contains(stdout.String(), "Generate of \"/dir/to/CHANGELOG.tpl\"")
}
