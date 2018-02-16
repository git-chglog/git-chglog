package chglog

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

type tagReader struct {
	client gitcmd.Client
	format string
	reTag  *regexp.Regexp
}

func newTagReader(client gitcmd.Client) *tagReader {
	return &tagReader{
		client: client,
		reTag:  regexp.MustCompile("tag: ([\\w\\.\\-_]+),?"),
	}
}

func (r *tagReader) ReadAll() ([]*Tag, error) {
	out, err := r.client.Exec(
		"log",
		"--tags",
		"--simplify-by-decoration",
		"--pretty=%D\t%at",
	)

	tags := []*Tag{}

	if err != nil {
		return tags, err
	}

	lines := strings.Split(out, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		tokens := strings.Split(line, "\t")

		res := r.reTag.FindAllStringSubmatch(tokens[0], -1)
		if len(res) == 0 {
			continue
		}

		ts, err := strconv.Atoi(tokens[1])
		if err != nil {
			continue
		}

		tags = append(tags, &Tag{
			Name: res[0][1],
			Date: time.Unix(int64(ts), 0),
		})
	}

	total := len(tags)

	for i, tag := range tags {
		var (
			next *RelateTag
			prev *RelateTag
		)

		if i > 0 {
			next = &RelateTag{
				Name: tags[i-1].Name,
				Date: tags[i-1].Date,
			}
		}

		if i+1 < total {
			prev = &RelateTag{
				Name: tags[i+1].Name,
				Date: tags[i+1].Date,
			}
		}

		tag.Next = next
		tag.Previous = prev
	}

	return tags, nil
}
