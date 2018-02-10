package chglog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

const (
	testRepoRoot = ".tmp"
)

func TestMain(m *testing.M) {
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func setup(dir string, setupRepo func(gitcmd.Client)) {
	cwd, _ := os.Getwd()

	testDir := filepath.Join(testRepoRoot, dir)
	os.MkdirAll(testDir, os.ModePerm)
	os.Chdir(testDir)

	git := gitcmd.New(nil)
	git.Exec("init")
	git.Exec("config", "user.name", "test_user")
	git.Exec("config", "user.email", "test@example.com")

	setupRepo(git)

	os.Chdir(cwd)
}

func cleanup() {
	os.RemoveAll(testRepoRoot)
}

func TestGeneratorWithTypeScopeSubject(t *testing.T) {
	assert := assert.New(t)

	testName := "type_scope_subject"

	setup(testName, func(git gitcmd.Client) {
		git.Exec("commit", "-m", "--allow-empty", "chore(*): First commit")
		git.Exec("commit", "--allow-empty", "-m", "feat(core): Add foo bar")
		git.Exec("commit", "--allow-empty", "-m", "docs(readme): Update usage #123")

		git.Exec("tag", "1.0.0")
		git.Exec("commit", "--allow-empty", "-m", "feat(parser): New some super options #333")
		git.Exec("commit", "--allow-empty", "-m", "Merge pull request #999 from tsuyoshiwada/patch-1")
		git.Exec("commit", "--allow-empty", "-m", "Merge pull request #1000 from tsuyoshiwada/patch-1")
		git.Exec("commit", "--allow-empty", "-m", "Revert \"feat(core): Add foo bar @mention and issue #987\"")

		git.Exec("tag", "1.1.0")
		git.Exec("commit", "--allow-empty", "-m", "feat(context): Online breaking change\n\nBREAKING CHANGE: Online breaking change message.")
		git.Exec("commit", "--allow-empty", "-m", "feat(router): Muliple breaking change\n\nThis is body,\n\nBREAKING CHANGE:\nMultiple\nbreaking\nchange message.")

		git.Exec("tag", "2.0.0-beta.0")
		git.Exec("commit", "--allow-empty", "-m", "refactor(context): gofmt")
		git.Exec("commit", "--allow-empty", "-m", "fix(core): Fix commit\n\nThis is body message.")
	})

	gen := NewGenerator(&Config{
		Bin:      "git",
		Path:     filepath.Join(testRepoRoot, testName),
		Template: filepath.Join("fixtures", testName+".md"),
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
			RevertPattern: "^Revert\\s\"([\\s\\S]*)\"\\s*This reverts commit (\\w*)\\.",
			RevertPatternMaps: []string{
				"Subject",
				"Hash",
			},
			NoteKeywords: []string{
				"BREAKING CHANGE",
			},
		},
	})

	buf := &bytes.Buffer{}
	gen.Generate(buf, "")

	assert.Equal(`<a name="2.0.0-beta.0"></a>
## 2.0.0-beta.0 (2018-02-10)

### Features

* **context:** Online breaking change
* **router:** Muliple breaking change

### BREAKING CHANGE

Multiple
breaking
change message.

Online breaking change message.



<a name="1.1.0"></a>
## 1.1.0 (2018-02-10)

### Features

* **parser:** New some super options #333

### Pull Requests

* Merge pull request #1000 from tsuyoshiwada/patch-1
* Merge pull request #999 from tsuyoshiwada/patch-1


<a name="1.0.0"></a>
## 1.0.0 (2018-02-10)

### Features

* **core:** Add foo bar`, strings.TrimSpace(buf.String()))
}
