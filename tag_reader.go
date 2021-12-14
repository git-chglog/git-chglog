package chglog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	gitcmd "github.com/tsuyoshiwada/go-gitcmd"
)

type tagReader struct {
	client    gitcmd.Client
	separator string
	reFilter  *regexp.Regexp
	sortBy    string
}

func newTagReader(client gitcmd.Client, filterPattern string, sort string) *tagReader {
	return &tagReader{
		client:    client,
		separator: "@@__CHGLOG__@@",
		reFilter:  regexp.MustCompile(filterPattern),
		sortBy:    sort,
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
		return tags, fmt.Errorf("failed to get git-tag: %w", err)
	}

	lines := strings.Split(out, "\n")

	for _, line := range lines {
		tokens := strings.Split(line, r.separator)

		if len(tokens) != 4 {
			continue
		}

		name := r.parseRefname(tokens[0])
		subject := r.parseSubject(tokens[1])
		date, errOfParseDate := r.parseDate(tokens[2])
		if errOfParseDate != nil {
			t, errOfParseDate := r.parseDate(tokens[3])
			if errOfParseDate != nil {
				return nil, errOfParseDate
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

	tags, err = r.sortTags(tags)
	if err != nil {
		return nil, err
	}

	r.assignPreviousAndNextTag(tags)

	return tags, nil
}

func (r *tagReader) sortTags(tags []*Tag) ([]*Tag, error) {
	switch r.sortBy {
	case "date":
		err := r.assignOrderToTags(tags)
		if err != nil {
			return nil, err
		}
		r.sortTagsByDate(tags)
	case "semver":
		r.filterSemVerTags(&tags)
		r.sortTagsBySemver(tags)
	}
	return tags, nil
}

func (*tagReader) filterSemVerTags(tags *[]*Tag) {
	// filter out any non-semver tags
	for i, t := range *tags {
		// remove leading v, since its so
		// common.
		name := t.Name
		name = strings.TrimPrefix(name, "v")

		// attempt semver parse, if not successful
		// remove it from tags slice.
		if _, err := semver.NewVersion(name); err != nil {
			*tags = append((*tags)[:i], (*tags)[i+1:]...)
		}
	}
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

func (r *tagReader) assignOrderToTags(tags []*Tag) error {
	out, err := r.client.Exec("log", "--oneline", "--decorate=short")
	if err != nil {
		return fmt.Errorf("failed to get git-log: %w", err)
	}

	lines := strings.Split(out, "\n")
	re := regexp.MustCompile(`tag:\s([^\s]*)[,)]`)
	for i, line := range lines {
		result := re.FindAllStringSubmatch(line, -1)
		for _, match := range result {
			tagOfCommit := match[1]
			for _, tag := range tags {
				if tagOfCommit == tag.Name {
					tag.Order = i
					break
				}
			}
		}
	}
	return nil
}

func (*tagReader) sortTagsByDate(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Order < tags[j].Order
	})
}

func (*tagReader) sortTagsBySemver(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool {
		semver1 := strings.TrimPrefix(tags[i].Name, "v")
		semver2 := strings.TrimPrefix(tags[j].Name, "v")
		v1 := semver.New(semver1)
		v2 := semver.New(semver2)
		return v2.LessThan(*v1)
	})
}
