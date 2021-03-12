package chglog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

type tagReader struct {
	client    gitcmd.Client
	format    string
	separator string
	reFilter  *regexp.Regexp
}

func newTagReader(client gitcmd.Client, filterPattern string) *tagReader {
	return &tagReader{
		client:    client,
		separator: "@@__CHGLOG__@@",
		reFilter:  regexp.MustCompile(filterPattern),
	}
}

func (r *tagReader) ReadRange(query string) ([]*Tag, error) {
	// First get all of the tags and fetch the author / date
	out, err := r.client.Exec(
		"for-each-ref",
		"--format",
		"%(refname)"+r.separator+"%(subject)"+r.separator+"%(taggerdate)"+r.separator+"%(authordate)",
		"refs/tags",
	)

	allTags := make(map[string]Tag)
	desiredTags := []*Tag{}

	if err != nil {
		return desiredTags, fmt.Errorf("failed to get git-tag: %s", err.Error())
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

		allTags[name] = Tag{
			Name:    name,
			Subject: subject,
			Date:    date,
		}
	}

	// after fetching all of the tags in the repo, walk the tree specified by the query
	// i.e. `git log --pretty=%D --simplify-by-decoration RefA..RefB`

	walkRefs, err := r.client.Exec(
		"log",
		"--pretty=%D",
		"--simplify-by-decoration",
		query)

	lines2 := strings.Split(walkRefs, "\n")
	for _, line2 := range lines2 {
		tokens := strings.Split(line2, ",")

		for _, observedRef := range tokens {
			refTokens := strings.Split(observedRef, " ")
			if refTokens[0] == "tag:" {
				desiredTags = append(desiredTags, &Tag{
					Date:    allTags[refTokens[1]].Date,
					Name:    refTokens[1],
					Subject: allTags[refTokens[1]].Subject,
				})
			}
		}

	}
	r.assignPreviousAndNextTag(desiredTags)
	return desiredTags, nil

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

		if r.reFilter != nil {
			if !r.reFilter.MatchString(name) {
				continue
			}
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
