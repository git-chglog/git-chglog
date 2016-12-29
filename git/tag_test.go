package git

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCurrentTag(t *testing.T) {
	assert := assert.New(t)
	tag, err := CurrentTag()
	assert.NoError(err)
	assert.NotEmpty(tag)
}

func TestPreviousTag(t *testing.T) {
	assert := assert.New(t)
	tag, err := PreviousTag("v0.0.1")
	assert.NoError(err)
	assert.NotEmpty(tag)
}