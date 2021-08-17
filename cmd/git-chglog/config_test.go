package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigNormalize(t *testing.T) {
	assert := assert.New(t)

	// basic
	config := &Config{
		Info: Info{
			RepositoryURL: "https://example.com/foo/bar/",
		},
	}

	err := config.Normalize(&CLIContext{
		ConfigPath: filepath.FromSlash("/test/config.yml"),
	})

	assert.Nil(err)
	assert.Equal("git", config.Bin)
	assert.Equal("https://example.com/foo/bar", config.Info.RepositoryURL)
	assert.Equal("/test/CHANGELOG.tpl.md", filepath.ToSlash(config.Template))

	// abs template
	cwd, _ := os.Getwd()

	config = &Config{
		Template: filepath.Join(cwd, "CHANGELOG.tpl.md"),
	}

	err = config.Normalize(&CLIContext{
		ConfigPath: filepath.Join(cwd, "test", "config.yml"),
	})

	assert.Nil(err)
	assert.Equal(filepath.Join(cwd, "CHANGELOG.tpl.md"), config.Template)
}

func TestConfigConvert(t *testing.T) {
	var patternInFile = "pattern in config"
	var patternInArgs = "pattern in cli"
	var titleInFile = "title in file"
	var titleInArgs = "title in args"
	assert := assert.New(t)
	// basic
	config := &Config{
		Info: Info{
			RepositoryURL: "https://example.com/foo/bar/",
			Title:         titleInFile,
		},
		Options: Options{TagFilterPattern: patternInFile},
	}
	cli := &CLIContext{TagFilterPattern: patternInArgs, Title: titleInArgs}
	cfg := config.Convert(cli)
	assert.Equal(cfg.Options.TagFilterPattern, patternInArgs)
	assert.Equal(cfg.Info.Title, titleInArgs)

	config = &Config{
		Info: Info{
			RepositoryURL: "https://example.com/foo/bar/",
			Title:         titleInFile,
		},
		Options: Options{TagFilterPattern: patternInFile},
	}
	cli = &CLIContext{}
	cfg = config.Convert(cli)
	assert.Equal(cfg.Options.TagFilterPattern, patternInFile)
	assert.Equal(cfg.Info.Title, titleInFile)

}
