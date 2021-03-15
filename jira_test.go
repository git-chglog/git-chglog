package chglog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJira(t *testing.T) {
	assert := assert.New(t)

	config := &Config{
		Options: &Options{
			Processor:                   nil,
			NextTag:                     "",
			TagFilterPattern:            "",
			CommitFilters:               nil,
			CommitSortBy:                "",
			CommitGroupBy:               "",
			CommitGroupSortBy:           "",
			CommitGroupTitleMaps:        nil,
			HeaderPattern:               "",
			HeaderPatternMaps:           nil,
			IssuePrefix:                 nil,
			RefActions:                  nil,
			MergePattern:                "",
			MergePatternMaps:            nil,
			RevertPattern:               "",
			RevertPatternMaps:           nil,
			NoteKeywords:                nil,
			JiraUsername:                "uuu",
			JiraToken:                   "ppp",
			JiraURL:                     "http://jira.com",
			JiraTypeMaps:                nil,
			JiraIssueDescriptionPattern: "",
		},
	}

	jira := NewJiraClient(config)
	issue, err := jira.GetJiraIssue("fake")
	assert.Nil(issue)
	assert.Error(err)
}
