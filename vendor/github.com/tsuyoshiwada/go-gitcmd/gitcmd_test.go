package gitcmd

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestGitClient(t *testing.T) {
	git := New(nil)

	if err := git.CanExec(); err != nil {
		t.Errorf("CanExec() should be return nil: %v\n", err)
	}

	if err := git.InsideWorkTree(); err != nil {
		t.Errorf("InsideWorkTree() should be return nil: %v\n", err)
	}

	out, err := git.Exec("rev-parse", "--git-dir")
	expected := ".git"

	if err != nil {
		t.Errorf("Exec() should be return nil: %v\n", err)
	}

	if out != expected {
		t.Errorf("got %v\nwant %v\n", out, expected)
	}
}

func TestAwareBin(t *testing.T) {
	bytes, _ := exec.Command("which", "git").Output()
	bin := strings.TrimSpace(string(bytes))

	git := New(&Config{
		Bin: bin,
	})

	if err := git.CanExec(); err != nil {
		t.Errorf("CanExec() should be return nil: %v\n", err)
	}
}

func TestNotfoundBin(t *testing.T) {
	git := New(&Config{
		Bin: "/notfound/git/bin",
	})

	if err := git.CanExec(); err == nil {
		t.Errorf("CanExec() should be return error")
	}
}

// Mocks
type MockClient struct {
	Client
	ReturnCanExec        func() error
	ReturnExec           func(string, ...string) (string, error)
	ReturnInsideWorkTree func() error
}

func (m *MockClient) CanExec() error {
	return m.ReturnCanExec()
}

func (m *MockClient) Exec(subcmd string, args ...string) (string, error) {
	return m.ReturnExec(subcmd, args...)
}

func (m *MockClient) InsideWorkTree() error {
	return m.ReturnInsideWorkTree()
}

func TestMock(t *testing.T) {
	var (
		errCanExec        = errors.New("CanExec() error")
		errExec           = errors.New("Exec() error")
		errInsideWorkTree = errors.New("InsideWorkTree() error")
	)

	git := &MockClient{}

	git.ReturnCanExec = func() error {
		return errCanExec
	}

	git.ReturnExec = func(subcmd string, args ...string) (string, error) {
		if subcmd == "log" {
			return "", errExec
		}
		return "", nil
	}

	git.ReturnInsideWorkTree = func() error {
		return errInsideWorkTree
	}

	if err := git.CanExec(); err == nil {
		t.Errorf("MockClient.CanExec() should be return error")
	}

	if _, err := git.Exec("log"); err == nil {
		t.Errorf("MockClient.Exec(\"log\") should be return error")
	}

	if err := git.InsideWorkTree(); err == nil {
		t.Errorf("MockClient.InsideWorkTree() should be return error")
	}
}
