package main

import (
	"fmt"
	"strings"
)

// ConfigBuilder ...
type ConfigBuilder interface {
	Builder
}

type configBuilderImpl struct{}

// NewConfigBuilder ...
func NewConfigBuilder() ConfigBuilder {
	return &configBuilderImpl{}
}

// Build ...
func (*configBuilderImpl) Build(ans *Answer) (string, error) {
	var msgFormat *CommitMessageFormat

	for _, ff := range formats {
		f, _ := ff.(*CommitMessageFormat)
		if f.display == ans.CommitMessageFormat {
			msgFormat = f
			break
		}
	}

	if msgFormat == nil {
		return "", fmt.Errorf("\"%s\" is an invalid commit message format", ans.CommitMessageFormat)
	}

	repoURL := strings.TrimRight(ans.RepositoryURL, "/")
	if repoURL == "" {
		repoURL = "\"\""
	}

	config := fmt.Sprintf(`style: %s
template: %s
info:
  title: CHANGELOG
  repository_url: %s
options:
  commits:
    # filters:
    #   Type:
    #     - feat
    #     - fix
    #     - perf
    #     - refactor
  commit_groups:
    # title_maps:
    #   feat: Features
    #   fix: Bug Fixes
    #   perf: Performance Improvements
    #   refactor: Code Refactoring
  header:
    pattern: "%s"
    pattern_maps:%s
  notes:
    keywords:
      - BREAKING CHANGE`,
		ans.Style,
		defaultTemplateFilename,
		repoURL,
		msgFormat.pattern,
		msgFormat.PatternMapString(),
	)

	return config, nil
}
