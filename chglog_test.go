package chglog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

var (
	cwd                string
	testRepoRoot       = ".tmp"
	internalTimeFormat = "2006-01-02 15:04:05"
)

type commitFunc = func(date, subject, body string)
type tagFunc = func(name string)

func TestMain(m *testing.M) {
	cwd, _ = os.Getwd()
	cleanup()
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func setup(dir string, setupRepo func(commitFunc, tagFunc, gitcmd.Client)) {
	testDir := filepath.Join(cwd, testRepoRoot, dir)

	os.RemoveAll(testDir)
	os.MkdirAll(testDir, os.ModePerm)
	os.Chdir(testDir)

	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	git := gitcmd.New(nil)
	git.Exec("init")
	git.Exec("config", "user.name", "test_user")
	git.Exec("config", "user.email", "test@example.com")

	var commit = func(date, subject, body string) {
		msg := subject
		if body != "" {
			msg += "\n\n" + body
		}
		t, _ := time.Parse(internalTimeFormat, date)
		d := t.Format("Mon Jan 2 15:04:05 2006 +0000")
		git.Exec("commit", "--allow-empty", "--date", d, "-m", msg)
	}

	var tag = func(name string) {
		git.Exec("tag", name)
	}

	setupRepo(commit, tag, git)

	os.Chdir(cwd)
}

func cleanup() {
	os.Chdir(cwd)
	os.RemoveAll(filepath.Join(cwd, testRepoRoot))
}

func TestGeneratorNotFoundTags(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, _ tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(*): New feature", "")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{},
	})

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "")
	assert.Error(err)
	assert.Contains(err.Error(), "git-tag does not exist")
	assert.Equal("", buf.String())
}

func TestGeneratorNotFoundCommits(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(*): New feature", "")
		tag("1.0.0")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{},
	})

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "foo")
	assert.Error(err)
	assert.Equal("", buf.String())
}

func TestGeneratorNotFoundCommitsOne(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "chore(*): First commit", "")
		tag("1.0.0")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{
			CommitFilters:        map[string][]string{},
			CommitSortBy:         "Scope",
			CommitGroupBy:        "Type",
			CommitGroupSortBy:    "Title",
			CommitGroupTitleMaps: map[string]string{},
			HeaderPattern:        "^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$",
			HeaderPatternMaps: []string{
				"Type",
				"Scope",
				"Subject",
			},
			IssuePrefix: []string{
				"#",
				"gh-",
			},
			RefActions:   []string{},
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

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "foo")
	assert.Error(err)
	assert.Contains(err.Error(), "\"foo\" was not found")
	assert.Equal("", buf.String())
}

func TestGeneratorWithTypeScopeSubject(t *testing.T) {
	assert := assert.New(t)
	testName := "type_scope_subject"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "chore(*): First commit", "")
		commit("2018-01-01 00:01:00", "feat(core): Add foo bar", "")
		commit("2018-01-01 00:02:00", "docs(readme): Update usage #123", "")
		tag("1.0.0")

		commit("2018-01-02 00:00:00", "feat(parser): New some super options #333", "")
		commit("2018-01-02 00:01:00", "Merge pull request #999 from tsuyoshiwada/patch-1", "")
		commit("2018-01-02 00:02:00", "Merge pull request #1000 from tsuyoshiwada/patch-1", "")
		commit("2018-01-02 00:03:00", "Revert \"feat(core): Add foo bar @mention and issue #987\"", "")
		tag("1.1.0")

		commit("2018-01-03 00:00:00", "feat(context): Online breaking change", "BREAKING CHANGE: Online breaking change message.")
		commit("2018-01-03 00:01:00", "feat(router): Muliple breaking change", `This is body,

BREAKING CHANGE:
Multiple
breaking
change message.`)
		tag("2.0.0-beta.0")

		commit("2018-01-04 00:00:00", "refactor(context): gofmt", "")
		commit("2018-01-04 00:01:00", "fix(core): Fix commit\n\nThis is body message.", "")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			Title:         "CHANGELOG Example",
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{
			CommitFilters: map[string][]string{
				"Type": []string{
					"feat",
					"fix",
				},
			},
			CommitSortBy:      "Scope",
			CommitGroupBy:     "Type",
			CommitGroupSortBy: "Title",
			CommitGroupTitleMaps: map[string]string{
				"feat": "Features",
				"fix":  "Bug Fixes",
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
			RefActions:   []string{},
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

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "")
	output := strings.Replace(strings.TrimSpace(buf.String()), "\r\n", "\n", -1)

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]

### Bug Fixes
- **core:** Fix commit


<a name="2.0.0-beta.0"></a>
## [2.0.0-beta.0] - 2018-01-03
### Features
- **context:** Online breaking change
- **router:** Muliple breaking change

### BREAKING CHANGE

Multiple
breaking
change message.

Online breaking change message.


<a name="1.1.0"></a>
## [1.1.0] - 2018-01-02
### Features
- **parser:** New some super options #333

### Reverts
- feat(core): Add foo bar @mention and issue #987

### Pull Requests
- Merge pull request #1000 from tsuyoshiwada/patch-1
- Merge pull request #999 from tsuyoshiwada/patch-1


<a name="1.0.0"></a>
## 1.0.0 - 2018-01-01
### Features
- **core:** Add foo bar


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/2.0.0-beta.0...HEAD
[2.0.0-beta.0]: https://github.com/git-chglog/git-chglog/compare/1.1.0...2.0.0-beta.0
[1.1.0]: https://github.com/git-chglog/git-chglog/compare/1.0.0...1.1.0`, output)
}

func TestGeneratorWithNextTag(t *testing.T) {
	assert := assert.New(t)
	testName := "type_scope_subject"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(core): version 1.0.0", "")
		tag("1.0.0")

		commit("2018-02-01 00:00:00", "feat(core): version 2.0.0", "")
		tag("2.0.0")

		commit("2018-03-01 00:00:00", "feat(core): version 3.0.0", "")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			Title:         "CHANGELOG Example",
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{
			NextTag: "3.0.0",
			CommitFilters: map[string][]string{
				"Type": []string{
					"feat",
				},
			},
			CommitSortBy:      "Scope",
			CommitGroupBy:     "Type",
			CommitGroupSortBy: "Title",
			CommitGroupTitleMaps: map[string]string{
				"feat": "Features",
			},
			HeaderPattern: "^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$",
			HeaderPatternMaps: []string{
				"Type",
				"Scope",
				"Subject",
			},
		},
	})

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "")
	output := strings.Replace(strings.TrimSpace(buf.String()), "\r\n", "\n", -1)

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="3.0.0"></a>
## [3.0.0] - 2018-03-01
### Features
- **core:** version 3.0.0


<a name="2.0.0"></a>
## [2.0.0] - 2018-02-01
### Features
- **core:** version 2.0.0


<a name="1.0.0"></a>
## 1.0.0 - 2018-01-01
### Features
- **core:** version 1.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/3.0.0...HEAD
[3.0.0]: https://github.com/git-chglog/git-chglog/compare/2.0.0...3.0.0
[2.0.0]: https://github.com/git-chglog/git-chglog/compare/1.0.0...2.0.0`, output)

	buf = &bytes.Buffer{}
	err = gen.Generate(buf, "3.0.0")
	output = strings.Replace(strings.TrimSpace(buf.String()), "\r\n", "\n", -1)

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="3.0.0"></a>
## [3.0.0] - 2018-03-01
### Features
- **core:** version 3.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/3.0.0...HEAD
[3.0.0]: https://github.com/git-chglog/git-chglog/compare/2.0.0...3.0.0`, output)
}

func TestGeneratorWithTagFiler(t *testing.T) {
	assert := assert.New(t)
	testName := "type_scope_subject"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(core): version dev-1.0.0", "")
		tag("dev-1.0.0")

		commit("2018-02-01 00:00:00", "feat(core): version v1.0.0", "")
		tag("v1.0.0")
	})

	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: filepath.Join(testRepoRoot, testName),
		Template:   filepath.Join(cwd, "testdata", testName+".md"),
		Info: &Info{
			Title:         "CHANGELOG Example",
			RepositoryURL: "https://github.com/git-chglog/git-chglog",
		},
		Options: &Options{
			TagFilterPattern: "^v",
			CommitFilters: map[string][]string{
				"Type": []string{
					"feat",
				},
			},
			CommitSortBy:      "Scope",
			CommitGroupBy:     "Type",
			CommitGroupSortBy: "Title",
			CommitGroupTitleMaps: map[string]string{
				"feat": "Features",
			},
			HeaderPattern: "^(\\w*)(?:\\(([\\w\\$\\.\\-\\*\\s]*)\\))?\\:\\s(.*)$",
			HeaderPatternMaps: []string{
				"Type",
				"Scope",
				"Subject",
			},
		},
	})

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "")

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.0"></a>
## v1.0.0 - 2018-02-01
### Features
- **core:** version v1.0.0
- **core:** version dev-1.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/v1.0.0...HEAD`, strings.TrimSpace(buf.String()))

}
