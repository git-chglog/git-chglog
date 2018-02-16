package chglog

import (
	"errors"
	"strings"
	"testing"
	"time"

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

	actual, err := newTagReader(client).ReadAll()
	assert.Nil(err)

	assert.Equal(
		[]*Tag{
			&Tag{
				Name: "v5.2.0-beta.1",
				Date: time.Unix(1518023112, 0),
				Next: nil,
				Previous: &RelateTag{
					Name: "2.0.0",
					Date: time.Unix(1517875200, 0),
				},
			},
			&Tag{
				Name: "2.0.0",
				Date: time.Unix(1517875200, 0),
				Next: &RelateTag{
					Name: "v5.2.0-beta.1",
					Date: time.Unix(1518023112, 0),
				},
				Previous: &RelateTag{
					Name: "v2.0.4-rc.1",
					Date: time.Unix(1517788800, 0),
				},
			},
			&Tag{
				Name: "v2.0.4-rc.1",
				Date: time.Unix(1517788800, 0),
				Next: &RelateTag{
					Name: "2.0.0",
					Date: time.Unix(1517875200, 0),
				},
				Previous: &RelateTag{
					Name: "2.0.4-beta.1",
					Date: time.Unix(1517702400, 0),
				},
			},
			&Tag{
				Name: "2.0.4-beta.1",
				Date: time.Unix(1517702400, 0),
				Next: &RelateTag{
					Name: "v2.0.4-rc.1",
					Date: time.Unix(1517788800, 0),
				},
				Previous: &RelateTag{
					Name: "hoge_fuga",
					Date: time.Unix(1517616000, 0),
				},
			},
			&Tag{
				Name: "hoge_fuga",
				Date: time.Unix(1517616000, 0),
				Next: &RelateTag{
					Name: "2.0.4-beta.1",
					Date: time.Unix(1517702400, 0),
				},
				Previous: &RelateTag{
					Name: "1.9.29-alpha.0",
					Date: time.Unix(1517529600, 0),
				},
			},
			&Tag{
				Name: "1.9.29-alpha.0",
				Date: time.Unix(1517529600, 0),
				Next: &RelateTag{
					Name: "hoge_fuga",
					Date: time.Unix(1517616000, 0),
				},
				Previous: nil,
			},
		},
		actual,
	)
}
