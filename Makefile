.PHONY: bootstrap
bootstrap: clean deps

.PHONY: deps
deps:
	dep ensure -v

.PHONY: clean
clean:
	rm -rf ./vendor/
	rm -rf ./git-chglog
	rm -rf $(GOPATH)/bin/git-chglog

.PHONY: bulid
build:
	go build -i -o git-chglog

.PHONY: test
test:
	go test -v `go list ./... | grep -v /vendor/`

.PHONY: install
install:
	go install ./cmd/git-chglog

.PHONY: chglog
chglog:
	git-chglog -c ./.chglog/config.yml
