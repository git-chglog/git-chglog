FROM golang:1.12.5-alpine

RUN apk add --no-cache git && go get -u github.com/git-chglog/git-chglog/cmd/git-chglog

ENTRYPOINT [ "git-chglog", "--version" ]
