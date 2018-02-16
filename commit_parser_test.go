package chglog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommitParserParse(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)

	mock := &mockClient{
		ReturnExec: func(subcmd string, args ...string) (string, error) {
			if subcmd != "log" {
				return "", errors.New("")
			}

			bytes, _ := ioutil.ReadFile(filepath.Join("testdata", "gitlog.txt"))

			return string(bytes), nil
		},
	}

	parser := newCommitParser(mock, &Config{
		Options: &Options{
			CommitFilters: map[string][]string{
				"Type": []string{
					"feat",
					"fix",
					"perf",
					"refactor",
				},
			},
			HeaderPattern: "^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$",
			HeaderPatternMaps: []string{
				"Type",
				"Scope",
				"Subject",
			},
			IssuePrefix: []string{
				"#",
				"gh-",
			},
			RefActions: []string{
				"close",
				"closes",
				"closed",
				"fix",
				"fixes",
				"fixed",
				"resolve",
				"resolves",
				"resolved",
			},
			MergePattern: "^Merge pull request #(\\d+) from (.*)$",
			MergePatternMaps: []string{
				"Ref",
				"Source",
			},
			RevertPattern: "^Revert \"([\\s\\S]*)\"$",
			RevertPatternMaps: []string{
				"Header",
			},
			NoteKeywords: []string{
				"BREAKING CHANGE",
			},
		},
	})

	commits, err := parser.Parse("HEAD")
	assert.Nil(err)
	assert.Equal([]*Commit{
		&Commit{
			Hash: &Hash{
				Long:  "65cf1add9735dcc4810dda3312b0792236c97c4e",
				Short: "65cf1add",
			},
			Author: &Author{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1514808000), 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1514808000), 0),
			},
			Merge:  nil,
			Revert: nil,
			Refs: []*Ref{
				&Ref{
					Action: "",
					Ref:    "123",
					Source: "",
				},
			},
			Notes:    []*Note{},
			Mentions: []string{},
			Header:   "feat(*): Add new feature #123",
			Type:     "feat",
			Scope:    "*",
			Subject:  "Add new feature #123",
			Body:     "",
		},
		&Commit{
			Hash: &Hash{
				Long:  "14ef0b6d386c5432af9292eab3c8314fa3001bc7",
				Short: "14ef0b6d",
			},
			Author: &Author{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1515153600), 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1515153600), 0),
			},
			Merge: &Merge{
				Ref:    "3",
				Source: "username/branchname",
			},
			Revert: nil,
			Refs: []*Ref{
				&Ref{
					Action: "",
					Ref:    "3",
					Source: "",
				},
				&Ref{
					Action: "Fixes",
					Ref:    "3",
					Source: "",
				},
				&Ref{
					Action: "Closes",
					Ref:    "1",
					Source: "",
				},
			},
			Notes: []*Note{
				&Note{
					Title: "BREAKING CHANGE",
					Body:  "This is breaking point message.",
				},
			},
			Mentions: []string{},
			Header:   "Merge pull request #3 from username/branchname",
			Type:     "",
			Scope:    "",
			Subject:  "",
			Body: `This is body message.

Fixes #3

Closes #1

BREAKING CHANGE: This is breaking point message.`,
		},
		&Commit{
			Hash: &Hash{
				Long:  "809a8280ffd0dadb0f4e7ba9fc835e63c37d6af6",
				Short: "809a8280",
			},
			Author: &Author{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517486400), 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517486400), 0),
			},
			Merge:  nil,
			Revert: nil,
			Refs:   []*Ref{},
			Notes:  []*Note{},
			Mentions: []string{
				"tsuyoshiwada",
				"hogefuga",
				"FooBarBaz",
			},
			Header:  "fix(controller): Fix cors configure",
			Type:    "fix",
			Scope:   "controller",
			Subject: "Fix cors configure",
			Body: `Has mention body

@tsuyoshiwada
@hogefuga
@FooBarBaz`,
		},
		&Commit{
			Hash: &Hash{
				Long:  "74824d6bd1470b901ec7123d13a76a1b8938d8d0",
				Short: "74824d6b",
			},
			Author: &Author{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517488587), 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517488587), 0),
			},
			Merge:  nil,
			Revert: nil,
			Refs: []*Ref{
				&Ref{
					Action: "Fixes",
					Ref:    "123",
					Source: "",
				},
				&Ref{
					Action: "Closes",
					Ref:    "456",
					Source: "username/repository",
				},
			},
			Notes: []*Note{
				&Note{
					Title: "BREAKING CHANGE",
					Body: fmt.Sprintf(`This is multiline breaking change note.
It is treated as the body of the Note until a mention or reference appears.

We also allow blank lines :)

Example:

%sjavascript
import { Controller } from 'hoge-fuga';

@autobind
class MyController extends Controller {
  constructor() {
    super();
  }
}
%s`, "```", "```"),
				},
			},
			Mentions: []string{},
			Header:   "fix(model): Remove hoge attributes",
			Type:     "fix",
			Scope:    "model",
			Subject:  "Remove hoge attributes",
			Body: fmt.Sprintf(`This mixed body message.

BREAKING CHANGE:
This is multiline breaking change note.
It is treated as the body of the Note until a mention or reference appears.

We also allow blank lines :)

Example:

%sjavascript
import { Controller } from 'hoge-fuga';

@autobind
class MyController extends Controller {
  constructor() {
    super();
  }
}
%s

Fixes #123
Closes username/repository#456`, "```", "```"),
		},
		&Commit{
			Hash: &Hash{
				Long:  "123456789735dcc4810dda3312b0792236c97c4e",
				Short: "12345678",
			},
			Author: &Author{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517488587), 0),
			},
			Committer: &Committer{
				Name:  "tsuyoshi wada",
				Email: "mail@example.com",
				Date:  time.Unix(int64(1517488587), 0),
			},
			Merge: nil,
			Revert: &Revert{
				Header: "fix(core): commit message",
			},
			Refs:     []*Ref{},
			Notes:    []*Note{},
			Mentions: []string{},
			Header:   "Revert \"fix(core): commit message\"",
			Type:     "",
			Scope:    "",
			Subject:  "",
			Body:     "This reverts commit f755db78dcdf461dc42e709b3ab728ceba353d1d.",
		},
	}, commits)
}
