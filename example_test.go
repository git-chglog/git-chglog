package chglog

import (
	"bytes"
	"fmt"
	"log"
)

func Example() {
	gen := NewGenerator(&Config{
		Bin:        "git",
		WorkingDir: ".",
		Template:   "CHANGELOG.tpl.md",
		Info: &Info{
			Title:         "CHANGELOG",
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

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(buf.String())
}
