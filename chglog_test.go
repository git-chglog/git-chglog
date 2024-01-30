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

	_ = os.RemoveAll(testDir)
	_ = os.MkdirAll(testDir, os.ModePerm)
	_ = os.Chdir(testDir)

	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	git := gitcmd.New(nil)
	_, _ = git.Exec("init")
	_, _ = git.Exec("config", "user.name", "test_user")
	_, _ = git.Exec("config", "user.email", "test@example.com")

	var commit = func(date, subject, body string) {
		msg := subject
		if body != "" {
			msg += "\n\n" + body
		}
		t, _ := time.Parse(internalTimeFormat, date)
		d := t.Format("Mon Jan 2 15:04:05 2006 +0000")
		_, _ = git.Exec("commit", "--allow-empty", "--date", d, "-m", msg)
	}

	var tag = func(name string) {
		_, _ = git.Exec("tag", name)
	}

	setupRepo(commit, tag, git)

	_ = os.Chdir(cwd)
}

func cleanup() {
	_ = os.Chdir(cwd)
	_ = os.RemoveAll(filepath.Join(cwd, testRepoRoot))
}

func TestGeneratorNotFoundTags(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, _ tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(*): New feature", "")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
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
	expected := strings.TrimSpace(buf.String())

	assert.Error(err)
	assert.Contains(err.Error(), "git-tag does not exist")
	assert.Equal("", expected)
}

func TestGeneratorNotFoundCommits(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(*): New feature", "")
		tag("1.0.0")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
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
	expected := strings.TrimSpace(buf.String())

	assert.Error(err)
	assert.Equal("", expected)
}

func TestGeneratorNotFoundCommitsOne(t *testing.T) {
	assert := assert.New(t)
	testName := "not_found"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "chore(*): First commit", "")
		tag("1.0.0")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
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
	expected := strings.TrimSpace(buf.String())

	assert.Error(err)
	assert.Contains(err.Error(), "\"foo\" was not found")
	assert.Equal("", expected)
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
		commit("2018-01-03 00:01:00", "feat(router): Multiple breaking change", `This is body,

BREAKING CHANGE:
Multiple
breaking
change message.`)
		tag("2.0.0-beta.0")

		commit("2018-01-04 00:00:00", "refactor(context): gofmt", "")
		commit("2018-01-04 00:01:00", "fix(core): Fix commit\n\nThis is body message.", "")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
			Bin:        "git",
			WorkingDir: filepath.Join(testRepoRoot, testName),
			Template:   filepath.Join(cwd, "testdata", testName+".md"),
			Info: &Info{
				Title:         "CHANGELOG Example",
				RepositoryURL: "https://github.com/git-chglog/git-chglog",
			},
			Options: &Options{
				Sort: "date",
				CommitFilters: map[string][]string{
					"Type": {
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
	expected := strings.TrimSpace(buf.String())

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]

### Bug Fixes
- **core:** Fix commit


<a name="2.0.0-beta.0"></a>
## [2.0.0-beta.0] - 2018-01-03
### Features
- **context:** Online breaking change
- **router:** Multiple breaking change

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
[1.1.0]: https://github.com/git-chglog/git-chglog/compare/1.0.0...1.1.0`, expected)
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

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
			Bin:        "git",
			WorkingDir: filepath.Join(testRepoRoot, testName),
			Template:   filepath.Join(cwd, "testdata", testName+".md"),
			Info: &Info{
				Title:         "CHANGELOG Example",
				RepositoryURL: "https://github.com/git-chglog/git-chglog",
			},
			Options: &Options{
				Sort:    "date",
				NextTag: "3.0.0",
				CommitFilters: map[string][]string{
					"Type": {
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
	expected := strings.TrimSpace(buf.String())

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
[2.0.0]: https://github.com/git-chglog/git-chglog/compare/1.0.0...2.0.0`, expected)

	buf = &bytes.Buffer{}
	err = gen.Generate(buf, "3.0.0")
	expected = strings.TrimSpace(buf.String())

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="3.0.0"></a>
## [3.0.0] - 2018-03-01
### Features
- **core:** version 3.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/3.0.0...HEAD
[3.0.0]: https://github.com/git-chglog/git-chglog/compare/2.0.0...3.0.0`, expected)
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

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
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
					"Type": {
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
	expected := strings.TrimSpace(buf.String())

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="v1.0.0"></a>
## v1.0.0 - 2018-02-01
### Features
- **core:** version v1.0.0
- **core:** version dev-1.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/v1.0.0...HEAD`, expected)

}

func TestGeneratorWithTrimmedBody(t *testing.T) {
	assert := assert.New(t)
	testName := "trimmed_body"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat: single line commit", "")
		commit("2018-01-01 00:01:00", "feat: multi-line commit", `
More details about the change and why it went in.

BREAKING CHANGE:

When using .TrimmedBody Notes are not included and can only appear in the Notes section.

Signed-off-by: First Last <first.last@mail.com>

Co-authored-by: dependabot-preview[bot] <27856297+dependabot-preview[bot]@users.noreply.github.com>`)

		commit("2018-01-01 00:00:00", "feat: another single line commit", "")
		tag("1.0.0")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
			Bin:        "git",
			WorkingDir: filepath.Join(testRepoRoot, testName),
			Template:   filepath.Join(cwd, "testdata", testName+".md"),
			Info: &Info{
				Title:         "CHANGELOG Example",
				RepositoryURL: "https://github.com/git-chglog/git-chglog",
			},
			Options: &Options{
				CommitFilters: map[string][]string{
					"Type": {
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
				NoteKeywords: []string{
					"BREAKING CHANGE",
				},
			},
		})

	buf := &bytes.Buffer{}
	err := gen.Generate(buf, "")
	expected := strings.TrimSpace(buf.String())

	assert.Nil(err)
	assert.Equal(`<a name="unreleased"></a>
## [Unreleased]


<a name="1.0.0"></a>
## 1.0.0 - 2018-01-01
### Features
- another single line commit
- multi-line commit
  More details about the change and why it went in.
- single line commit

### BREAKING CHANGE

When using .TrimmedBody Notes are not included and can only appear in the Notes section.


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/1.0.0...HEAD`, expected)
}

func TestGeneratorWithSprig(t *testing.T) {
	assert := assert.New(t)
	testName := "with_sprig"

	setup(testName, func(commit commitFunc, tag tagFunc, _ gitcmd.Client) {
		commit("2018-01-01 00:00:00", "feat(core): version 1.0.0", "")
		tag("1.0.0")

		commit("2018-02-01 00:00:00", "feat(core): version 2.0.0", "")
		tag("2.0.0")

		commit("2018-03-01 00:00:00", "feat(core): version 3.0.0", "")
	})

	gen := NewGenerator(NewLogger(os.Stdout, os.Stderr, false, true),
		&Config{
			Bin:        "git",
			WorkingDir: filepath.Join(testRepoRoot, testName),
			Template:   filepath.Join(cwd, "testdata", testName+".md"),
			Info: &Info{
				Title:         "CHANGELOG Example",
				RepositoryURL: "https://github.com/git-chglog/git-chglog",
			},
			Options: &Options{
				Sort:    "date",
				NextTag: "3.0.0",
				CommitFilters: map[string][]string{
					"Type": {
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
	expected := strings.TrimSpace(buf.String())

	assert.Nil(err)
	assert.Equal(`My Changelog
<a name="unreleased"></a>
## [Unreleased]


<a name="3.0.0"></a>
## [3.0.0] - 2018-03-01
### Features
- **CORE:** version 3.0.0


<a name="2.0.0"></a>
## [2.0.0] - 2018-02-01
### Features
- **CORE:** version 2.0.0


<a name="1.0.0"></a>
## 1.0.0 - 2018-01-01
### Features
- **CORE:** version 1.0.0


[Unreleased]: https://github.com/git-chglog/git-chglog/compare/3.0.0...HEAD
[3.0.0]: https://github.com/git-chglog/git-chglog/compare/2.0.0...3.0.0
[2.0.0]: https://github.com/git-chglog/git-chglog/compare/1.0.0...2.0.0`, expected)

}

func TestUniqueOlderCommits(t *testing.T) {
	tests := []struct {
		name        string
		commits     []*Commit
		fields      []string
		expected    []*Commit
		expectError assert.ErrorAssertionFunc
	}{
		{
			name: "Duplication detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				// {Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			expectError: assert.NoError,
		},
		{
			name: "2 duplications detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log2"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				// {Hash: &Hash{Short: "1"}, Subject: "log1"},
				// {Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log2"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			expectError: assert.NoError,
		},
		{
			name: "A double duplication detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				// {Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				// {Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			expectError: assert.NoError,
		},
		{
			name: "No duplicates",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			expectError: assert.NoError,
		},
		{
			name:        "No commits (empty)",
			commits:     []*Commit{},
			fields:      []string{"Subject"},
			expected:    nil,
			expectError: assert.NoError,
		},
		{
			name:        "No commits (nil)",
			commits:     nil,
			fields:      []string{"Subject"},
			expected:    nil,
			expectError: assert.NoError,
		},
		{
			name: "Empty fields (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      []string{},
			expected:    nil,
			expectError: assert.Error,
		},
		{
			name: "nil for fields (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      nil,
			expected:    nil,
			expectError: assert.Error,
		},
		{
			name: "bad field (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      []string{"NoSuchField"},
			expected:    nil,
			expectError: assert.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := uniqueOlderCommits(
				test.commits,
				test.fields...,
			)
			test.expectError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestUniqueNewerCommits(t *testing.T) {
	tests := []struct {
		name        string
		commits     []*Commit
		fields      []string
		expected    []*Commit
		expectError assert.ErrorAssertionFunc
	}{
		{
			name: "Duplication detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				// {Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			expectError: assert.NoError,
		},
		{
			name: "2 duplications detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log2"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				// {Hash: &Hash{Short: "3"}, Subject: "log2"}, // Duplicated!
				// {Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			expectError: assert.NoError,
		},
		{
			name: "A double duplication detected",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				// {Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				// {Hash: &Hash{Short: "4"}, Subject: "log1"}, // Duplicated!
			},
			expectError: assert.NoError,
		},
		{
			name: "No duplicates",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields: []string{"Subject"},
			expected: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			expectError: assert.NoError,
		},
		{
			name:        "No commits (empty)",
			commits:     []*Commit{},
			fields:      []string{"Subject"},
			expected:    nil,
			expectError: assert.NoError,
		},
		{
			name:        "No commits (nil)",
			commits:     nil,
			fields:      []string{"Subject"},
			expected:    nil,
			expectError: assert.NoError,
		},
		{
			name: "Empty fields (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      []string{},
			expected:    nil,
			expectError: assert.Error,
		},
		{
			name: "nil for fields (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      nil,
			expected:    nil,
			expectError: assert.Error,
		},
		{
			name: "bad field (error)",
			commits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log1"}, // Duplicated!
				{Hash: &Hash{Short: "4"}, Subject: "log4"},
			},
			fields:      []string{"NoSuchField"},
			expected:    nil,
			expectError: assert.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := uniqueNewerCommits(
				test.commits,
				test.fields...,
			)
			test.expectError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestHasDuplicatedCommit(t *testing.T) {
	tests := []struct {
		name          string
		targetCommits []*Commit
		commit        *Commit
		fields        []string
		expected      bool
		expectError   assert.ErrorAssertionFunc
	}{
		{
			name: "Duplication found",
			targetCommits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
			},
			commit:      &Commit{Hash: &Hash{Short: "99"}, Subject: "log2"},
			fields:      []string{"Subject"},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name: "Duplication not found",
			targetCommits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
			},
			commit:      &Commit{Hash: &Hash{Short: "99"}, Subject: "log99"},
			fields:      []string{"Subject"},
			expected:    false,
			expectError: assert.NoError,
		},
		{
			name:          "No commits (empty)",
			targetCommits: []*Commit{},
			commit:        &Commit{Hash: &Hash{Short: "99"}, Subject: "log1"},
			fields:        []string{"Subject"},
			expected:      false,
			expectError:   assert.NoError,
		},
		{
			name:          "No commits (nil)",
			targetCommits: []*Commit{},
			commit:        &Commit{Hash: &Hash{Short: "99"}, Subject: "log1"},
			fields:        []string{"Subject"},
			expected:      false,
			expectError:   assert.NoError,
		},
		{
			name: "Empty fields: Nothing to compare and everything is duplicated",
			targetCommits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
			},
			commit:      &Commit{Hash: &Hash{Short: "99"}, Subject: "log1"},
			fields:      []string{},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name: "nil for fields: Nothing to compare and everything is duplicated",
			targetCommits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
			},
			commit:      &Commit{Hash: &Hash{Short: "99"}, Subject: "log1"},
			fields:      nil,
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name: "bad field (error)",
			targetCommits: []*Commit{
				{Hash: &Hash{Short: "1"}, Subject: "log1"},
				{Hash: &Hash{Short: "2"}, Subject: "log2"},
				{Hash: &Hash{Short: "3"}, Subject: "log3"},
			},
			commit:      &Commit{Hash: &Hash{Short: "99"}, Subject: "log1"},
			fields:      []string{"NoSuchField"},
			expected:    false,
			expectError: assert.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := hasDuplicatedCommit(
				test.targetCommits,
				test.commit,
				test.fields,
			)
			test.expectError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestIsDuplicatedCommit(t *testing.T) {
	tests := []struct {
		name        string
		commit1     *Commit
		commit2     *Commit
		fields      []string
		expected    bool
		expectError assert.ErrorAssertionFunc
	}{
		{
			name:        "Duplicated for single field",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log1"},
			fields:      []string{"Subject"},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name:        "Not duplicted for single field",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2"},
			fields:      []string{"Subject"},
			expected:    false,
			expectError: assert.NoError,
		},
		{
			name:        "Duplicated for multiple fields",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Scope: "chore", Subject: "log1", JiraIssueID: "FOO-1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Scope: "chore", Subject: "log1", JiraIssueID: "FOO-1"},
			fields:      []string{"Scope", "Subject", "JiraIssueID"},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name:        "Not duplicted for multiple fields",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Scope: "chore", Subject: "log1", JiraIssueID: "FOO-1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Scope: "chore", Subject: "log1", JiraIssueID: "FOO-2"},
			fields:      []string{"Scope", "Subject", "JiraIssueID"},
			expected:    false,
			expectError: assert.NoError,
		},
		{
			name:        "Nested",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1", Author: &Author{Name: "user1"}},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2", Author: &Author{Name: "user1"}},
			fields:      []string{"Author.Name"},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name:        "Nested: dereference nil",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1", Author: &Author{Name: "user1"}},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2", Author: nil},
			fields:      []string{"Author.Name"},
			expected:    false,
			expectError: assert.NoError,
		},
		{
			name:        "Empty fields: Nothing to compare and everything is duplicated",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2"},
			fields:      []string{},
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name:        "nil for fields: Nothing to compare and everything is duplicated",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2"},
			fields:      nil,
			expected:    true,
			expectError: assert.NoError,
		},
		{
			name:        "bad field (error)",
			commit1:     &Commit{Hash: &Hash{Short: "1"}, Subject: "log1"},
			commit2:     &Commit{Hash: &Hash{Short: "2"}, Subject: "log2"},
			fields:      []string{"NoSuchField"},
			expected:    false,
			expectError: assert.Error,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := isDuplicatedCommit(
				test.commit1,
				test.commit2,
				test.fields,
			)
			test.expectError(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
