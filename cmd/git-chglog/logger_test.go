package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

func TestLoggerLogSilent(t *testing.T) {
	color.NoColor = false
	assert := assert.New(t)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	logger := NewLogger(stdout, stderr, true, false)
	logger.Log(":+1:Hello, World! :)")
	assert.Equal("", stdout.String())
}

func TestLoggerLog(t *testing.T) {
	color.NoColor = false
	assert := assert.New(t)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	logger := NewLogger(stdout, stderr, false, false)
	logger.Log(":+1:Hello, World! :)")
	assert.Equal(emoji.Sprint(":+1:Hello, World! :)\n"), stdout.String())
}

func TestLoggerLogNoEmoji(t *testing.T) {
	color.NoColor = false
	assert := assert.New(t)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	logger := NewLogger(stdout, stderr, false, true)
	logger.Log(":+1:Hello, World! :)")
	assert.Equal(fmt.Sprint("Hello, World! :)\n"), stdout.String())
}

func TestLoggerError(t *testing.T) {
	color.NoColor = false
	assert := assert.New(t)

	prefix := color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()

	// Basic
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	logger := NewLogger(stdout, stderr, false, false)
	logger.Error("This is error message!! :dog:")
	assert.Equal("", stdout.String())
	assert.Equal(emoji.Sprint(fmt.Sprintf("%s %s\n", prefix(" ERROR "), color.RedString("This is error message!! :dog:"))), stderr.String())

	// Silent
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	logger = NewLogger(stdout, stderr, true, false)
	logger.Error("Foo")
	assert.Equal("", stdout.String())
	assert.NotEqual("", stderr.String())

	// NoEmoji
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	logger = NewLogger(stdout, stderr, true, true)
	logger.Error("HOGE :hand:")
	assert.Equal("", stdout.String())
	assert.NotContains(stderr.String(), emoji.Sprint(":hand:"))
}
