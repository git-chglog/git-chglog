package main

import (
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
		ConfigPath: "/test/config.yml",
	})

	assert.Nil(err)
	assert.Equal("git", config.Bin)
	assert.Equal("https://example.com/foo/bar", config.Info.RepositoryURL)
	assert.Equal("/test/CHANGELOG.tpl.md", config.Template)

	// abs template
	config = &Config{
		Template: "/CHANGELOG.tpl.md",
	}

	err = config.Normalize(&CLIContext{
		ConfigPath: "/test/config.yml",
	})

	assert.Nil(err)
	assert.Equal("/CHANGELOG.tpl.md", config.Template)
}
