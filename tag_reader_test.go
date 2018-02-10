package chglog

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagReader(t *testing.T) {
	assert := assert.New(t)
	client := &mockClient{
		ReturnExec: func(subcmd string, args ...string) (string, error) {
			if subcmd != "log" {
				return "", errors.New("")
			}
			return strings.Join([]string{
				"",
				"tag: v5.2.0-beta.1, origin/labs/router\t1518023112",
				"tag: 2.0.0\t1517875200",
				"tag: v2.0.4-rc.1\t1517788800",
				"tag: 2.0.4-beta.1\t1517702400",
				"tag: hoge_fuga\t1517616000",
				"tag: 1.9.29-alpha.0\t1517529600",
				"hoge\t0",
				"foo\t0",
			}, "\n"), nil
		},
	}

	res, err := newTagReader(client).ReadAll()
	assert.Nil(err)

	actual := make([]string, len(res))
	for i, tag := range res {
		actual[i] = tag.Name
	}

	assert.Equal([]string{
		"v5.2.0-beta.1",
		"2.0.0",
		"v2.0.4-rc.1",
		"2.0.4-beta.1",
		"hoge_fuga",
		"1.9.29-alpha.0",
	}, actual)
}
