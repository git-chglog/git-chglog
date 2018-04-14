package main

import (
	"fmt"
	"net/url"

	chglog "github.com/git-chglog/git-chglog"
)

// ProcessorFactory ...
type ProcessorFactory struct {
	hostRegistry map[string]string
}

// NewProcessorFactory ...
func NewProcessorFactory() *ProcessorFactory {
	return &ProcessorFactory{
		hostRegistry: map[string]string{
			"github":    "github.com",
			"gitlab":    "gitlab.com",
			"bitbucket": "bitbucket.org",
		},
	}
}

// Create ...
func (factory *ProcessorFactory) Create(config *Config) (chglog.Processor, error) {
	obj, err := url.Parse(config.Info.RepositoryURL)
	if err != nil {
		return nil, err
	}

	host := obj.Host

	if config.Style != "" {
		if styleHost, ok := factory.hostRegistry[config.Style]; ok {
			host = styleHost
		}
	}

	switch host {
	case "github.com":
		return &chglog.GitHubProcessor{
			Host: fmt.Sprintf("%s://%s", obj.Scheme, obj.Host),
		}, nil
	case "gitlab.com":
		return &chglog.GitLabProcessor{
			Host: fmt.Sprintf("%s://%s", obj.Scheme, obj.Host),
		}, nil
	case "bitbucket.org":
		return &chglog.BitbucketProcessor{
			Host: fmt.Sprintf("%s://%s", obj.Scheme, obj.Host),
		}, nil
	default:
		return nil, nil
	}
}
