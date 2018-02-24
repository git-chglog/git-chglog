package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitMessageFormatPatternMaps(t *testing.T) {
	assert := assert.New(t)

	f := &CommitMessageFormat{
		PatternMaps: []string{
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
		PatternMaps: []string{},
	}

	assert.Equal(" []", f.PatternMapString())
}
