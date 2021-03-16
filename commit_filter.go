package chglog

import (
	"strings"
)

func commitFilter(commits []*Commit, filters map[string][]string, noCaseSensitive bool) []*Commit { //nolint:gocyclo
	// NOTE(khos2ow): this function is over our cyclomatic complexity goal.
	// Be wary when adding branches, and look for functionality that could
	// be reasonably moved into an injected dependency.

	res := []*Commit{}

	for _, commit := range commits {
		include := false

		if len(filters) == 0 {
			include = true
		}

		for key, values := range filters {
			prop, ok := dotGet(commit, key)
			if !ok {
				include = false
				break
			}

			str, ok := prop.(string)
			if !ok {
				include = false
				break
			}

			if noCaseSensitive {
				str = strings.ToLower(str)
			}

			exist := false

			for _, val := range values {
				if noCaseSensitive {
					val = strings.ToLower(val)
				}

				if str == val {
					exist = true
				}
			}

			if !exist {
				include = false
				break
			}

			include = true
		}

		if include {
			res = append(res, commit)
		}
	}

	return res
}
