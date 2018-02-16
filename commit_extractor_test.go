package chglog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitExtractor(t *testing.T) {
	assert := assert.New(t)

	extractor := newCommitExtractor(&Options{
		CommitSortBy:      "Scope",
		CommitGroupBy:     "Type",
		CommitGroupSortBy: "Title",
		CommitGroupTitleMaps: map[string]string{
			"bar": "BAR",
		},
	})

	fixtures := []*Commit{
		// [0]
		&Commit{
			Type:   "foo",
			Scope:  "c",
			Header: "1",
			Notes:  []*Note{},
		},
		// [1]
		&Commit{
			Type:   "foo",
			Scope:  "b",
			Header: "2",
			Notes: []*Note{
				{"note1-title", "note1-body"},
				{"note2-title", "note2-body"},
			},
		},
		// [2]
		&Commit{
			Type:   "bar",
			Scope:  "d",
			Header: "3",
			Notes: []*Note{
				{"note1-title", "note1-body"},
				{"note3-title", "note3-body"},
			},
		},
		// [3]
		&Commit{
			Type:   "foo",
			Scope:  "a",
			Header: "4",
			Notes: []*Note{
				{"note4-title", "note4-body"},
			},
		},
		// [4]
		&Commit{
			Type:   "",
			Scope:  "",
			Header: "Merge1",
			Notes:  []*Note{},
			Merge: &Merge{
				Ref:    "123",
				Source: "merges/merge1",
			},
		},
		// [5]
		&Commit{
			Type:   "",
			Scope:  "",
			Header: "Revert1",
			Notes:  []*Note{},
			Revert: &Revert{
				Header: "REVERT1",
			},
		},
	}

	commitGroups, mergeCommits, revertCommits, noteGroups := extractor.Extract(fixtures)

	assert.Equal([]*CommitGroup{
		&CommitGroup{
			RawTitle: "bar",
			Title:    "BAR",
			Commits: []*Commit{
				fixtures[2],
			},
		},
		&CommitGroup{
			RawTitle: "foo",
			Title:    "Foo",
			Commits: []*Commit{
				fixtures[3],
				fixtures[1],
				fixtures[0],
			},
		},
	}, commitGroups)

	assert.Equal([]*Commit{
		fixtures[4],
	}, mergeCommits)

	assert.Equal([]*Commit{
		fixtures[5],
	}, revertCommits)

	assert.Equal([]*NoteGroup{
		&NoteGroup{
			Title: "note1-title",
			Notes: []*Note{
				fixtures[1].Notes[0],
				fixtures[2].Notes[0],
			},
		},
		&NoteGroup{
			Title: "note2-title",
			Notes: []*Note{
				fixtures[1].Notes[1],
			},
		},
		&NoteGroup{
			Title: "note3-title",
			Notes: []*Note{
				fixtures[2].Notes[1],
			},
		},
		&NoteGroup{
			Title: "note4-title",
			Notes: []*Note{
				fixtures[3].Notes[0],
			},
		},
	}, noteGroups)
}
