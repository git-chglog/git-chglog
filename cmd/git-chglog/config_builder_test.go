package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigBulider(t *testing.T) {
	assert := assert.New(t)
	builder := NewConfigBuilder()

	out, err := builder.Build(&Answer{
		RepositoryURL:       "https://github.com/git-chglog/git-chglog/git-chglog/",
		Style:               styleNone,
		CommitMessageFormat: fmtGitBasic.Display,
		Template:            tplStandard,
	})

	assert.Nil(err)
	assert.Contains(out, "style: none")
	assert.Contains(out, "template: CHANGELOG.tpl.md")
	assert.Contains(out, "  repository_url: https://github.com/git-chglog/git-chglog/git-chglog")
	assert.Contains(out, fmt.Sprintf("    pattern: \"%s\"", fmtGitBasic.Pattern))
	assert.Contains(out, fmt.Sprintf(
		`    pattern_maps:
      - %s
      - %s`,
		fmtGitBasic.PatternMaps[0],
		fmtGitBasic.PatternMaps[1],
	))
}

func TestConfigBuliderEmptyRepoURL(t *testing.T) {
	assert := assert.New(t)
	builder := NewConfigBuilder()

	out, err := builder.Build(&Answer{
		RepositoryURL:       "",
		Style:               styleNone,
		CommitMessageFormat: fmtGitBasic.Display,
		Template:            tplStandard,
	})

	assert.Nil(err)
	assert.Contains(out, "  repository_url: \"\"")
}

func TestConfigBuliderInvalidFormat(t *testing.T) {
	assert := assert.New(t)
	builder := NewConfigBuilder()

	_, err := builder.Build(&Answer{
		RepositoryURL:       "",
		Style:               styleNone,
		CommitMessageFormat: "",
		Template:            tplStandard,
	})

	assert.Contains(err.Error(), "invalid commit message format")
}
