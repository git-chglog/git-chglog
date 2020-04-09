package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitMessageFormatPatternMaps(t *testing.T) {
	assert := assert.New(t)

	f := &CommitMessageFormat{
		patternMaps: []string{
			"Type",
			"Scope",
			"Subject",
		},
	}

	assert.Equal(`
      - Type
      - Scope
      - Subject`, f.PatternMapString())

	f = &CommitMessageFormat{
		patternMaps: []string{},
	}

	assert.Equal(" []", f.PatternMapString())
}

func TestCommitMessageFormatFilterTypes(t *testing.T) {
	assert := assert.New(t)

	f := &CommitMessageFormat{
		typeSamples: []typeSample{
			{"feat", "Features"}, {"fix", "Bug Fixes"},
			{"perf", "Performance Improvements"}, {"refactor", "Code Refactoring"},
		},
	}

	assert.Equal(`
    #     - feat
    #     - fix
    #     - perf
    #     - refactor`, f.FilterTypesString())

	f = &CommitMessageFormat{
		patternMaps: []string{},
	}

	assert.Equal(" []", f.FilterTypesString())
}

func TestCommitMessageFormatTitleMaps(t *testing.T) {
	assert := assert.New(t)

	f := &CommitMessageFormat{
		typeSamples: []typeSample{
			{"feat", "Features"}, {"fix", "Bug Fixes"},
			{"perf", "Performance Improvements"}, {"refactor", "Code Refactoring"},
		},
	}

	assert.Equal(`
    #   feat: Features
    #   fix: Bug Fixes
    #   perf: Performance Improvements
    #   refactor: Code Refactoring`, f.TitleMapsString())

	f = &CommitMessageFormat{
		patternMaps: []string{},
	}

	assert.Equal(" []", f.TitleMapsString())
}
