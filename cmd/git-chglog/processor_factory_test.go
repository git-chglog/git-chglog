package main

import (
	"testing"

	chglog "github.com/git-chglog/git-chglog"
	"github.com/stretchr/testify/assert"
)

func TestProcessorFactory(t *testing.T) {
	assert := assert.New(t)
	factory := NewProcessorFactory()

	processor, err := factory.Create(&Config{
		Info: Info{
			RepositoryURL: "https://example.com/owner/repo",
		},
	})

	assert.Nil(err)
	assert.Nil(processor)
}

func TestProcessorFactoryForGitHub(t *testing.T) {
	assert := assert.New(t)
	factory := NewProcessorFactory()

	// github.com
	processor, err := factory.Create(&Config{
		Info: Info{
			RepositoryURL: "https://github.com/owner/repo",
		},
	})

	assert.Nil(err)
	assert.Equal(
		&chglog.GitHubProcessor{
			Host: "https://github.com",
		},
		processor,
	)

	// ghe
	processor, err = factory.Create(&Config{
		Style: "github",
		Info: Info{
			RepositoryURL: "https://ghe-example.com/owner/repo",
		},
	})

	assert.Nil(err)
	assert.Equal(
		&chglog.GitHubProcessor{
			Host: "https://ghe-example.com",
		},
		processor,
	)
}
