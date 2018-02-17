package chglog

import (
	"strings"
)

type tagSelector struct{}

func newTagSelector() *tagSelector {
	return &tagSelector{}
}

func (s *tagSelector) Select(tags []*Tag, query string) ([]*Tag, string, error) {
	tokens := strings.Split(query, "..")

	switch len(tokens) {
	case 1:
		return s.selectSingleTag(tags, tokens[0])
	case 2:
		old := tokens[0]
		new := tokens[1]
		if old == "" && new == "" {
			return nil, "", nil
		} else if old == "" {
			return s.selectBeforeTags(tags, new)
		} else if new == "" {
			return s.selectAfterTags(tags, old)
		}
		return s.selectRangeTags(tags, tokens[0], tokens[1])
	}

	return nil, "", errFailedQueryParse
}

func (s *tagSelector) selectSingleTag(tags []*Tag, token string) ([]*Tag, string, error) {
	var from string

	for i, tag := range tags {
		if tag.Name == token {
			if i+1 < len(tags) {
				from = tags[i+1].Name
			}
			return []*Tag{tag}, from, nil
		}
	}

	return nil, "", nil
}

func (*tagSelector) selectBeforeTags(tags []*Tag, token string) ([]*Tag, string, error) {
	var (
		res    []*Tag
		from   string
		enable bool
	)

	for i, tag := range tags {
		if tag.Name == token {
			enable = true
		}

		if enable {
			res = append(res, tag)
			from = ""
			if i+1 < len(tags) {
				from = tags[i+1].Name
			}
		}
	}

	if len(res) == 0 {
		return res, "", errNotFoundTag
	}

	return res, from, nil
}

func (*tagSelector) selectAfterTags(tags []*Tag, token string) ([]*Tag, string, error) {
	var (
		res  []*Tag
		from string
	)

	for i, tag := range tags {
		res = append(res, tag)
		from = ""
		if i+1 < len(tags) {
			from = tags[i+1].Name
		}

		if tag.Name == token {
			break
		}
	}

	if len(res) == 0 {
		return res, "", errNotFoundTag
	}

	return res, from, nil
}

func (s *tagSelector) selectRangeTags(tags []*Tag, old string, new string) ([]*Tag, string, error) {
	var (
		res    []*Tag
		from   string
		enable bool
	)

	for i, tag := range tags {
		if tag.Name == new {
			enable = true
		}

		if enable {
			from = ""
			if i+1 < len(tags) {
				from = tags[i+1].Name
			}
			res = append(res, tag)
		}

		if tag.Name == old {
			enable = false
		}
	}

	if len(res) == 0 {
		return res, "", errNotFoundTag
	}

	return res, from, nil
}
