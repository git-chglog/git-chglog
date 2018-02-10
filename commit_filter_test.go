package chglog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitFilter(t *testing.T) {
	assert := assert.New(t)

	pickCommitSubjects := func(arr []*Commit) []string {
		res := make([]string, len(arr))
		for i, commit := range arr {
			res[i] = commit.Subject
		}
		return res
	}

	fixtures := []*Commit{
		&Commit{
			Type:    "foo",
			Scope:   "hoge",
			Subject: "1",
		},
		&Commit{
			Type:    "foo",
			Scope:   "fuga",
			Subject: "2",
		},
		&Commit{
			Type:    "bar",
			Scope:   "hoge",
			Subject: "3",
		},
		&Commit{
			Type:    "bar",
			Scope:   "fuga",
			Subject: "4",
		},
	}

	assert.Equal(
		[]string{
			"1",
			"2",
			"3",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{})),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
			"3",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"foo", "bar"},
		})),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"foo"},
		})),
	)

	assert.Equal(
		[]string{
			"3",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"bar"},
		})),
	)

	assert.Equal(
		[]string{
			"2",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Scope": {"fuga"},
		})),
	)

	assert.Equal(
		[]string{
			"3",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type":  {"bar"},
			"Scope": {"hoge"},
		})),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type":  {"foo"},
			"Scope": {"fuga", "hoge"},
		})),
	)
}
