package main

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/fatih/color"
	emoji "gopkg.in/kyokomi/emoji.v1"
)

// Logger ...
type Logger struct {
	stdout  io.Writer
	stderr  io.Writer
	silent  bool
	noEmoji bool
	reEmoji *regexp.Regexp
}

// NewLogger ...
func NewLogger(stdout, stderr io.Writer, silent, noEmoji bool) *Logger {
	return &Logger{
		stdout:  stdout,
		stderr:  stderr,
		silent:  silent,
		noEmoji: noEmoji,
		reEmoji: regexp.MustCompile(":[\\w\\+_\\-]+:\\s?"),
	}
}

// Log ...
func (l *Logger) Log(msg string) {
	if !l.silent {
		l.log(l.stdout, msg+"\n")
	}
}

// Error ...
func (l *Logger) Error(msg string) {
	prefix := color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
	l.log(l.stderr, fmt.Sprintf("%s %s\n", prefix(" ERROR "), color.RedString(msg)))
}

func (l *Logger) log(w io.Writer, msg string) {
	var printer func(io.Writer, ...interface{}) (int, error)

	if l.noEmoji {
		msg = l.reEmoji.ReplaceAllString(msg, "")
		printer = fmt.Fprint
	} else {
		printer = emoji.Fprint
	}

	if _, err := printer(w, msg); err != nil {
		log.Fatalln(err)
	}
}
