NAME    := k8stail
VERSION := $(shell git tag | sort -V -r | head -n1)-next
COMMIT  := $(shell git rev-parse HEAD)
DATE    := $(shell date "+%Y-%m-%dT%H:%M:%S%z")

SRCS    := $(shell find . -name '*.go' -type f)
LDFLAGS := -ldflags="-s -w -X \"main.version=$(VERSION)\" -X \"main.commit=$(COMMIT)\" -X \"main.date=$(DATE)\""

.DEFAULT_GOAL := bin/$(NAME)

export GO111MODULE=on

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	go test -coverpkg=./... -coverprofile=coverage.txt -v ./...

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -v
