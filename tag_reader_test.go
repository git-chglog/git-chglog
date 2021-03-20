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
			if subcmd != "for-each-ref" {
				return "", errors.New("")
			}
			return strings.Join([]string{
				"",
				"refs/tags/v2.0.4-beta.1@@__CHGLOG__@@Release v2.0.4-beta.1@@__CHGLOG__@@Thu Feb 1 00:00:00 2018 +0000@@__CHGLOG__@@",
				"refs/tags/4.4.3@@__CHGLOG__@@This is tag subject@@__CHGLOG__@@@@__CHGLOG__@@Fri Feb 2 00:00:00 2018 +0000",
				"refs/tags/4.4.4@@__CHGLOG__@@Release 4.4.4@@__CHGLOG__@@Fri Feb 2 10:00:40 2018 +0000@@__CHGLOG__@@",
				"refs/tags/v2.0.4-beta.2@@__CHGLOG__@@Release v2.0.4-beta.2@@__CHGLOG__@@Sat Feb 3 12:15:00 2018 +0000@@__CHGLOG__@@",
				"refs/tags/5.0.0-rc.0@@__CHGLOG__@@Release 5.0.0-rc.0@@__CHGLOG__@@Sat Feb 3 12:30:10 2018 +0000@@__CHGLOG__@@",
				"refs/tags/hoge_fuga@@__CHGLOG__@@Invalid semver tag name@@__CHGLOG__@@Mon Mar 12 12:30:10 2018 +0000@@__CHGLOG__@@",
				"hoge@@__CHGLOG__@@",
			}, "\n"), nil
		},
	}

	actual, err := newTagReader(client, "", "date").ReadAll()
	assert.Nil(err)

	assert.Equal(
		[]*Tag{
			{
				Name:    "hoge_fuga",
				Subject: "Invalid semver tag name",
				Date:    time.Date(2018, 3, 12, 12, 30, 10, 0, time.UTC),
				Next:    nil,
				Previous: &RelateTag{
					Name:    "5.0.0-rc.0",
					Subject: "Release 5.0.0-rc.0",
					Date:    time.Date(2018, 2, 3, 12, 30, 10, 0, time.UTC),
				},
			},
			{
				Name:    "5.0.0-rc.0",
				Subject: "Release 5.0.0-rc.0",
				Date:    time.Date(2018, 2, 3, 12, 30, 10, 0, time.UTC),
				Next: &RelateTag{
					Name:    "hoge_fuga",
					Subject: "Invalid semver tag name",
					Date:    time.Date(2018, 3, 12, 12, 30, 10, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "v2.0.4-beta.2",
					Subject: "Release v2.0.4-beta.2",
					Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				},
			},
			{
				Name:    "v2.0.4-beta.2",
				Subject: "Release v2.0.4-beta.2",
				Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "5.0.0-rc.0",
					Subject: "Release 5.0.0-rc.0",
					Date:    time.Date(2018, 2, 3, 12, 30, 10, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "4.4.4",
					Subject: "Release 4.4.4",
					Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				},
			},
			{
				Name:    "4.4.4",
				Subject: "Release 4.4.4",
				Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				Next: &RelateTag{
					Name:    "v2.0.4-beta.2",
					Subject: "Release v2.0.4-beta.2",
					Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "4.4.3",
					Subject: "This is tag subject",
					Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Name:    "4.4.3",
				Subject: "This is tag subject",
				Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "4.4.4",
					Subject: "Release 4.4.4",
					Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "v2.0.4-beta.1",
					Subject: "Release v2.0.4-beta.1",
					Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Name:    "v2.0.4-beta.1",
				Subject: "Release v2.0.4-beta.1",
				Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "4.4.3",
					Subject: "This is tag subject",
					Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				},
				Previous: nil,
			},
		},
		actual,
	)

	actual, err = newTagReader(client, "", "semver").ReadAll()
	assert.Nil(err)

	assert.Equal(
		[]*Tag{
			{
				Name:    "5.0.0-rc.0",
				Subject: "Release 5.0.0-rc.0",
				Date:    time.Date(2018, 2, 3, 12, 30, 10, 0, time.UTC),
				Next:    nil,
				Previous: &RelateTag{
					Name:    "4.4.4",
					Subject: "Release 4.4.4",
					Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				},
			},
			{
				Name:    "4.4.4",
				Subject: "Release 4.4.4",
				Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				Next: &RelateTag{
					Name:    "5.0.0-rc.0",
					Subject: "Release 5.0.0-rc.0",
					Date:    time.Date(2018, 2, 3, 12, 30, 10, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "4.4.3",
					Subject: "This is tag subject",
					Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Name:    "4.4.3",
				Subject: "This is tag subject",
				Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "4.4.4",
					Subject: "Release 4.4.4",
					Date:    time.Date(2018, 2, 2, 10, 0, 40, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "v2.0.4-beta.2",
					Subject: "Release v2.0.4-beta.2",
					Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				},
			},
			{
				Name:    "v2.0.4-beta.2",
				Subject: "Release v2.0.4-beta.2",
				Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "4.4.3",
					Subject: "This is tag subject",
					Date:    time.Date(2018, 2, 2, 0, 0, 0, 0, time.UTC),
				},
				Previous: &RelateTag{
					Name:    "v2.0.4-beta.1",
					Subject: "Release v2.0.4-beta.1",
					Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Name:    "v2.0.4-beta.1",
				Subject: "Release v2.0.4-beta.1",
				Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "v2.0.4-beta.2",
					Subject: "Release v2.0.4-beta.2",
					Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				},
				Previous: nil,
			},
		},
		actual,
	)

	actualFiltered, errFiltered := newTagReader(client, "^v", "date").ReadAll()
	assert.Nil(errFiltered)
	assert.Equal(
		[]*Tag{
			{
				Name:    "v2.0.4-beta.2",
				Subject: "Release v2.0.4-beta.2",
				Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				Next:    nil,
				Previous: &RelateTag{
					Name:    "v2.0.4-beta.1",
					Subject: "Release v2.0.4-beta.1",
					Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			{
				Name:    "v2.0.4-beta.1",
				Subject: "Release v2.0.4-beta.1",
				Date:    time.Date(2018, 2, 1, 0, 0, 0, 0, time.UTC),
				Next: &RelateTag{
					Name:    "v2.0.4-beta.2",
					Subject: "Release v2.0.4-beta.2",
					Date:    time.Date(2018, 2, 3, 12, 15, 0, 0, time.UTC),
				},
				Previous: nil,
			},
		},
		actualFiltered,
	)
}
