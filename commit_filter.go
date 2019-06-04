package chglog

func commitFilter(commits []*Commit, filters map[string][]string) []*Commit {
	res := []*Commit{}
	expandedCommits := []*Commit{}

	for _, commit := range commits {
		// expand squashed entries
		expandedCommits = append(expandedCommits, commit)
		if len(commit.AllHeaders) > 0 {
			expandedCommits = append(expandedCommits, commit.AllHeaders...)
		}
	}

	for _, commit := range expandedCommits {
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

			exist := false

			for _, val := range values {
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
