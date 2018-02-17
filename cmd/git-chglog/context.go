package main

import (
	"io"
)

// Context ...
type Context struct {
	WorkingDir string
	Stdout     io.Writer
	Stderr     io.Writer
	ConfigPath string
	OutputPath string
	Silent     bool
	NoColor    bool
	NoEmoji    bool
	Query      string
}
