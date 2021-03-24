VERSION?=$$(git describe --tags --always)

LDFLAGS="-X main.version=$(VERSION)"

.PHONY: clean
clean:
	rm -rf ./dist/
	rm -rf ./git-chglog
	rm -rf $(GOPATH)/bin/git-chglog
	rm -rf cover.out

.PHONY: build
build:
	go build -ldflags=$(LDFLAGS) -o git-chglog ./cmd/git-chglog

.PHONY: test
test:
	go test -covermode atomic -coverprofile cover.out -v ./...

.PHONY: install
install:
	go install ./cmd/git-chglog

.PHONY: changelog
changelog: build
	./git-chglog --next-tag $(VERSION) -o CHANGELOG.md

.PHONY: lint
lint:
	@golangci-lint run
