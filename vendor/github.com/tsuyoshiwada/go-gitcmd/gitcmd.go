// Package gitcmd is providing for tiny git command wrapper.
package gitcmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Config for git command
type Config struct {
	Bin string // default "git"
}

// Client of git command
type Client interface {
	CanExec() error
	Exec(string, ...string) (string, error)
	InsideWorkTree() error
}

type clientImpl struct {
	config *Config
}

// New git command client
func New(config *Config) Client {
	bin := "git"

	if config != nil {
		if config.Bin != "" {
			bin = config.Bin
		}
	}

	return &clientImpl{
		config: &Config{
			Bin: bin,
		},
	}
}

// CanExec check whether the git command is executable
func (client *clientImpl) CanExec() error {
	_, err := exec.LookPath(client.config.Bin)
	if err != nil {
		return fmt.Errorf("\"%s\" does not exists", client.config.Bin)
	}
	return nil
}

// Exec executes the git command
func (client *clientImpl) Exec(subcmd string, args ...string) (string, error) {
	arr := append([]string{subcmd}, args...)

	var out bytes.Buffer
	cmd := exec.Command(client.config.Bin, arr...)
	cmd.Stdout = &out
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() != 0 {
				return "", err
			}
		}
	}

	return strings.TrimRight(strings.TrimSpace(out.String()), "\000"), nil
}

// InsideWorkTree check whether the current working directory is inside the git repository
func (client *clientImpl) InsideWorkTree() error {
	out, err := client.Exec("rev-parse", "--is-inside-work-tree")
	if err != nil {
		return err
	}

	if out != "true" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return fmt.Errorf("\"%s\" is no git repository", cwd)
	}

	return nil
}
