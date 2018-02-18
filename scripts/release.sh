#!/bin/bash

if [[ "$TRAVIS_BRANCH" == "master" ]]; then
  gox -os="darwin linux windows" -arch="amd64 386" -output="$(pwd)/dist/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/git-chglog
  ghr --username git-chglog --token $GITHUB_TOKEN --replace `grep 'Version =' ./cmd/git-chglog/version.go | sed -E 's/.*"(.+)"$$/\1/'` dist/
fi
