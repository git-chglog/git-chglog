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
		&Commit{
			Type:    "Bar",
			Scope:   "hogera",
			Subject: "5",
		},
	}

	assert.Equal(
		[]string{
			"1",
			"2",
			"3",
			"4",
			"5",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{}, false)),
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
		}, false)),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
			"3",
			"4",
			"5",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"foo", "bar"},
		}, true)),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"foo"},
		}, false)),
	)

	assert.Equal(
		[]string{
			"3",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"bar"},
		}, false)),
	)

	assert.Equal(
		[]string{
			"3",
			"4",
			"5",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type": {"bar"},
		}, true)),
	)

	assert.Equal(
		[]string{
			"2",
			"4",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Scope": {"fuga"},
		}, false)),
	)

	assert.Equal(
		[]string{
			"3",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type":  {"bar"},
			"Scope": {"hoge"},
		}, false)),
	)

	assert.Equal(
		[]string{
			"1",
			"2",
		},
		pickCommitSubjects(commitFilter(fixtures, map[string][]string{
			"Type":  {"foo"},
			"Scope": {"fuga", "hoge"},
		}, false)),
	)
}
