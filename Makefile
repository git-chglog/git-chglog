.PHONY: bootstrap
bootstrap: clean

.PHONY: clean
clean:
	rm -rf ./dist/
	rm -rf ./git-chglog
	rm -rf $(GOPATH)/bin/git-chglog
	rm -rf cover.out

.PHONY: bulid
build:
	go build -o git-chglog github.com/git-chglog/git-chglog/cmd/git-chglog

.PHONY: test
test:
	go test -v ./...

.PHONY: coverage
coverage:
	goverage -coverprofile=cover.out `go list ./...`
	go tool cover -func=cover.out
	@rm -rf cover.out

.PHONY: install
install:
	go install ./cmd/git-chglog

.PHONY: changelog
changelog:
	@git-chglog --next-tag $(tag) $(tag)
