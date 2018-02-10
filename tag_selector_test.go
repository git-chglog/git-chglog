package chglog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagSelector(t *testing.T) {
	assert := assert.New(t)
	assert.True(true)
	selector := newTagSelector()

	fixtures := []*Tag{
		&Tag{Name: "2.2.12-rc.12"},
		&Tag{Name: "2.1.0"},
		&Tag{Name: "v2.0.0-beta.1"},
		&Tag{Name: "v1.2.9"},
		&Tag{Name: "v1.0.0"},
	}

	table := map[string][]string{
		// Single
		"2.2.12-rc.12": []string{
			"2.2.12-rc.12",
			"2.1.0",
		},
		"v2.0.0-beta.1": []string{
			"v2.0.0-beta.1",
			"v1.2.9",
		},
		"v1.0.0": []string{
			"v1.0.0",
			"",
		},
		// ~ <tag>
		"..2.1.0": []string{
			"2.1.0",
			"v2.0.0-beta.1",
			"v1.2.9",
			"v1.0.0",
			"",
		},
		"..v1.0.0": []string{
			"v1.0.0",
			"",
		},
		// <tag> ~
		"v2.0.0-beta.1..": []string{
			"2.2.12-rc.12",
			"2.1.0",
			"v2.0.0-beta.1",
			"v1.2.9",
		},
		"2.2.12-rc.12..": []string{
			"2.2.12-rc.12",
			"2.1.0",
		},
		"v1.0.0..": []string{
			"2.2.12-rc.12",
			"2.1.0",
			"v2.0.0-beta.1",
			"v1.2.9",
			"v1.0.0",
			"",
		},
		// <tag> ~ <tag>
		"v1.0.0..2.2.12-rc.12": []string{
			"2.2.12-rc.12",
			"2.1.0",
			"v2.0.0-beta.1",
			"v1.2.9",
			"v1.0.0",
			"",
		},
		"v1.0.0..v2.0.0-beta.1": []string{
			"v2.0.0-beta.1",
			"v1.2.9",
			"v1.0.0",
			"",
		},
		"v1.2.9..2.1.0": []string{
			"2.1.0",
			"v2.0.0-beta.1",
			"v1.2.9",
			"v1.0.0",
		},
	}

	for query, expected := range table {
		list, from, err := selector.Select(fixtures, query)
		actual := make([]string, len(list))
		for i, tag := range list {
			actual[i] = tag.Name
		}

		assert.Nil(err)
		assert.Equal(expected[0:len(expected)-1], actual)
		assert.Equal(expected[len(expected)-1], from)
	}
}
