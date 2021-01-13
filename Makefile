.PHONY: clean
clean:
	rm -rf ./dist/
	rm -rf ./git-chglog
	rm -rf $(GOPATH)/bin/git-chglog
	rm -rf cover.out

.PHONY: build
build:
	go build -i -o git-chglog

.PHONY: test
test:
	go test -coverprofile cover.out -v `go list ./...`

.PHONY: install
install:
	go install ./cmd/git-chglog

.PHONY: changelog
changelog:
	@git-chglog --next-tag $(tag) $(tag)
