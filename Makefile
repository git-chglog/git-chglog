# Build variables
VERSION ?= $(shell git describe --tags --always)

# Go variables
GO      ?= go
GOOS    ?= $(shell $(GO) env GOOS)
GOARCH  ?= $(shell $(GO) env GOARCH)
GOHOST  ?= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO)

LDFLAGS ?= "-X main.version=$(VERSION)"

# Docker variables
DEFAULT_TAG  ?= $(shell echo "$(VERSION)" | tr -d 'v')
DOCKER_IMAGE := quay.io/git-chglog/git-chglog
DOCKER_TAG   ?= $(DEFAULT_TAG)

.PHONY: all
all: help

###############
##@ Development

.PHONY: clean
clean:   ## Clean workspace
	@ $(MAKE) --no-print-directory log-$@
	rm -rf ./dist/
	rm -rf ./git-chglog
	rm -rf $(GOPATH)/bin/git-chglog
	rm -rf cover.out

.PHONY: test
test:   ## Run tests
	@ $(MAKE) --no-print-directory log-$@
	$(GOHOST) test -covermode atomic -coverprofile cover.out -v ./...

.PHONY: lint
lint:   ## Run linters
	@ $(MAKE) --no-print-directory log-$@
	golangci-lint run

#########
##@ Build

.PHONY: build
build:   ## Build git-chglog
	@ $(MAKE) --no-print-directory log-$@
	CGO_ENABLED=0 $(GOHOST) build -ldflags=$(LDFLAGS) -o git-chglog ./cmd/git-chglog

.PHONY: install
install:   ## Install git-chglog
	@ $(MAKE) --no-print-directory log-$@
	$(GOHOST) install ./cmd/git-chglog

.PHONY: docker
docker: build   ## Build Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker build --pull --tag $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: push
push:   ## Push Docker image
	@ $(MAKE) --no-print-directory log-$@
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

###########
##@ Release

.PHONY: changelog
changelog: build   ## Generate changelog
	@ $(MAKE) --no-print-directory log-$@
	./git-chglog --next-tag $(VERSION) -o CHANGELOG.md

.PHONY: release
release: changelog   ## Release a new tag
	@ $(MAKE) --no-print-directory log-$@
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
	git tag $(VERSION)
	git push origin master $(VERSION)

########
##@ Help

.PHONY: help
help:   ## Display this help
	@awk \
		-v "col=\033[36m" -v "nocol=\033[0m" \
		' \
			BEGIN { \
				FS = ":.*##" ; \
				printf "Usage:\n  make %s<target>%s\n", col, nocol \
			} \
			/^[a-zA-Z_-]+:.*?##/ { \
				printf "  %s%-12s%s %s\n", col, $$1, nocol, $$2 \
			} \
			/^##@/ { \
				printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol \
			} \
		' $(MAKEFILE_LIST)

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'
