BIN := "./bin/abf"
DOCKER_IMG="abf:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X 'github.com/razielsd/antibruteforce/app/cmd.version=develop' -X 'github.com/razielsd/antibruteforce/app/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X 'github.com/razielsd/antibruteforce/app/cmd.gitHash=$(GIT_HASH)'

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)"

run:
	$(BIN) server

build-img:
	docker build \
			--build-arg=LDFLAGS="$(LDFLAGS)" \
			-t $(DOCKER_IMG) \
			-f build/Dockerfile .

run-img:
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./app/...

lint:
	golangci-lint run ./app/...

.PHONY: build build-img run run-img version test lint