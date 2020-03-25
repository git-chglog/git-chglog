package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"testing"
)

var gAssert *assert.Assertions

func mock_app_action(c *cli.Context) error {
	assert := gAssert
	assert.Equal("c.yml", c.String("config"))
	assert.Equal("^v", c.String("tag-filter-pattern"))
	assert.Equal("o.md", c.String("output"))
	assert.Equal("v5", c.String("next-tag"))
	assert.True(c.Bool("silent"))
	assert.True(c.Bool("no-color"))
	assert.True(c.Bool("no-emoji"))
	return nil
}

func TestCreateApp(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
	gAssert = assert

	app := CreateApp(mock_app_action)
	args := []string{
		"git-chglog",
		"--silent",
		"--no-color",
		"--no-emoji",
		"--config", "c.yml",
		"--output", "o.md",
		"--next-tag", "v5",
		"--tag-filter-pattern", "^v",
	}
	app.Run(args)
}
