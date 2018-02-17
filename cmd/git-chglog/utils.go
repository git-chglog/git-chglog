package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var reSSH = regexp.MustCompile("^\\w+@([\\w\\.\\-]+):([\\w\\.\\-]+)\\/([\\w\\.\\-]+)$")

func remoteOriginURLToHTTP(rawurl string) string {
	if rawurl == "" {
		return ""
	}

	rawurl = strings.TrimSuffix(rawurl, ".git")

	// for normal url format
	originURL, err := url.Parse(rawurl)

	if err == nil {
		scheme := originURL.Scheme
		if scheme != "http" {
			scheme = "https"
		}
		return fmt.Sprintf(
			"%s://%s%s",
			scheme,
			originURL.Host,
			originURL.Path,
		)
	}

	// for `user@server:repo.git`
	res := reSSH.FindAllStringSubmatch(rawurl, -1)
	if len(res) > 0 {
		return fmt.Sprintf(
			"https://%s/%s/%s",
			res[0][1],
			res[0][2],
			res[0][3],
		)
	}

	return ""
}
