package chglog

import gitcmd "github.com/tsuyoshiwada/go-gitcmd"

type mockClient struct {
	gitcmd.Client
	ReturnExec func(string, ...string) (string, error)
}

func (m *mockClient) Exec(subcmd string, args ...string) (string, error) {
	return m.ReturnExec(subcmd, args...)
}
