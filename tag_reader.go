package chglog

import (
	"fmt"
	"sort"
	"strings"
	"time"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

type tagReader struct {
	client    gitcmd.Client
	format    string
	separator string
}

func newTagReader(client gitcmd.Client) *tagReader {
	return &tagReader{
		client:    client,
		separator: "@@__CHGLOG__@@",
	}
}

func (r *tagReader) ReadAll() ([]*Tag, error) {
	out, err := r.client.Exec(
		"for-each-ref",
		"--format",
		"%(refname)"+r.separator+"%(subject)"+r.separator+"%(taggerdate)"+r.separator+"%(authordate)",
		"refs/tags",
	)

	tags := []*Tag{}

	if err != nil {
		return tags, fmt.Errorf("failed to get git-tag: %s", err.Error())
	}

	lines := strings.Split(out, "\n")

	for _, line := range lines {
		tokens := strings.Split(line, r.separator)

		if len(tokens) != 4 {
			continue
		}

		name := r.parseRefname(tokens[0])
		subject := r.parseSubject(tokens[1])
		date, err := r.parseDate(tokens[2])
		if err != nil {
			t, err2 := r.parseDate(tokens[3])
			if err2 != nil {
				return nil, err2
			}
			date = t
		}

		tags = append(tags, &Tag{
			Name:    name,
			Subject: subject,
			Date:    date,
		})
	}

	r.sortTags(tags)
	r.assignPreviousAndNextTag(tags)

	return tags, nil
}

func (*tagReader) parseRefname(input string) string {
	return strings.Replace(input, "refs/tags/", "", 1)
}

func (*tagReader) parseSubject(input string) string {
	return strings.TrimSpace(input)
}

func (*tagReader) parseDate(input string) (time.Time, error) {
	return time.Parse("Mon Jan 2 15:04:05 2006 -0700", input)
}

func (*tagReader) assignPreviousAndNextTag(tags []*Tag) {
	total := len(tags)

	for i, tag := range tags {
		var (
			next *RelateTag
			prev *RelateTag
		)

		if i > 0 {
			next = &RelateTag{
				Name:    tags[i-1].Name,
				Subject: tags[i-1].Subject,
				Date:    tags[i-1].Date,
			}
		}

		if i+1 < total {
			prev = &RelateTag{
				Name:    tags[i+1].Name,
				Subject: tags[i+1].Subject,
				Date:    tags[i+1].Date,
			}
		}

		tag.Next = next
		tag.Previous = prev
	}
}

func (*tagReader) sortTags(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool {
		return !tags[i].Date.Before(tags[j].Date)
	})
}
