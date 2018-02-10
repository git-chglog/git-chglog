package chglog

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func workdir(dir string) func() {
	cwd, _ := os.Getwd()
	os.Chdir(dir)
	return func() {
		os.Chdir(cwd)
	}
}

func TestGenerator(t *testing.T) {
	t.Skip("TODO: test")
	assert := assert.New(t)
	assert.True(true)
}
