BIN := "./bin/abf"
DOCKER_IMG?="abf:develop"
ABF_PORT?=19187
CONTAINER_NAME="abf-container"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X 'github.com/razielsd/antibruteforce/app/cmd.version=develop' -X 'github.com/razielsd/antibruteforce/app/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X 'github.com/razielsd/antibruteforce/app/cmd.gitHash=$(GIT_HASH)'

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)"

run:
	$(BIN) server

build-img:
	docker build \
			--build-arg LDFLAGS="$(LDFLAGS)" \
			--build-arg ABF_PORT="$(ABF_PORT)" \
			-t $(DOCKER_IMG) \
			-f build/Dockerfile .

run-img:
	docker run --name=$(CONTAINER_NAME) -p $(ABF_PORT):$(ABF_PORT) --env ABF_ADDR="0.0.0.0:$(ABF_PORT)" $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test ./app/...

test-int:
	go test -race -tags integration -race ./app/...

test-int-coverage:
	go test -race -tags integration -race ./app/... `go list ./app/... | grep -v examples` -coverprofile=coverage.txt -covermode=atomic

test-coverage:
	go test -race ./app/... `go list ./app/... | grep -v examples` -coverprofile=coverage.txt -covermode=atomic

test100:
	go test -race -count 100 ./app/...

lint:
	golangci-lint run ./app/...

.PHONY: build build-img run run-img version test lint