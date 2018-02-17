package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteOriginURLToHTTP(t *testing.T) {
	assert := assert.New(t)

	table := [][]string{
		{"git://github.com/owner0/repo0.git", "https://github.com/owner0/repo0"},
		{"git@github.com:owner1/repo1.git", "https://github.com/owner1/repo1"},
		{"git+ssh://git@github.com/owner2/repo2.git", "https://github.com/owner2/repo2"},
		{"https://github.com/owner3/repo3.git", "https://github.com/owner3/repo3"},
		{"", ""},
	}

	for _, v := range table {
		assert.Equal(v[1], remoteOriginURLToHTTP(v[0]))
	}
}
