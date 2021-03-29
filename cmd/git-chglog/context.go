package main

import (
	"io"
)

// CLIContext ...
type CLIContext struct {
	WorkingDir       string
	Stdout           io.Writer
	Stderr           io.Writer
	ConfigPath       string
	Template         string
	RepositoryURL    string
	OutputPath       string
	Silent           bool
	NoColor          bool
	NoEmoji          bool
	NoCaseSensitive  bool
	Query            string
	NextTag          string
	TagFilterPattern string
	JiraUsername     string
	JiraToken        string
	JiraURL          string
	Paths            []string
	Sort             string
}

// InitContext ...
type InitContext struct {
	WorkingDir string
	Stdout     io.Writer
	Stderr     io.Writer
}
