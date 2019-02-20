version := $(shell git describe --exact-match --tags $(git log -n1 --pretty='%h') 2> /dev/null || echo 'latest')
vcs_ref := $(shell git rev-parse HEAD)

GO_FILES = $(shell find . -type f -name '*.go')

.PHONY: all
all: clean kroki

kroki: $(GO_FILES)
	go build -o $@ -ldflags "-s -w -X main.version=${version} -X main.commit=${vcs_ref}"

.PHONY: lint
lint: $(GO_FILES)
	golangci-lint run ./...

.PHONY: test
test: $(GO_FILES)
	go test -race ./...

.PHONY: clean
clean:
	go clean
	rm -rf kroki
