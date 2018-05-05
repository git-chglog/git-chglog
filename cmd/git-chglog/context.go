package main

import (
	"io"
)

// CLIContext ...
type CLIContext struct {
	WorkingDir string
	Stdout     io.Writer
	Stderr     io.Writer
	ConfigPath string
	OutputPath string
	Silent     bool
	NoColor    bool
	NoEmoji    bool
	Query      string
	NextTag    string
}

// InitContext ...
type InitContext struct {
	WorkingDir string
	Stdout     io.Writer
	Stderr     io.Writer
}
