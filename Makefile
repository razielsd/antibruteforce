BIN := "./bin/abf"
DOCKER_IMG?="abf:develop"
ABF_PORT?=19187
CONTAINER_NAME="abf-container"
DOCKER_IMG_TEST="abftest:latest"
CONTAINER_TEST="abftest"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X 'github.com/razielsd/antibruteforce/app/cmd.version=develop' -X 'github.com/razielsd/antibruteforce/app/cmd.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X 'github.com/razielsd/antibruteforce/app/cmd.gitHash=$(GIT_HASH)'
LOCAL_PATH := $(shell pwd)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)"

build-img:
	docker build \
			--build-arg LDFLAGS="$(LDFLAGS)" \
			--build-arg ABF_PORT="$(ABF_PORT)" \
			-t $(DOCKER_IMG) \
			-f build/Dockerfile .

build-test-img:
	docker build \
			-t $(DOCKER_IMG_TEST) \
			-f build/tests.Dockerfile .

run:
	$(BIN) server

run-img:
	docker run --rm --name=$(CONTAINER_NAME) -p $(ABF_PORT):$(ABF_PORT) --env ABF_ADDR="0.0.0.0:$(ABF_PORT)" $(DOCKER_IMG)

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

test-e2e:
	export ABF_BIN="$(LOCAL_PATH)/bin/abf" && export ABF_ADDR="127.0.0.1:49999" && go test -tags e2e ./tests/e2e/...

test-img: build-test-img
	docker run --rm -d --name=$(CONTAINER_TEST) abftest:latest
	docker exec $(CONTAINER_TEST) make test
	docker stop $(CONTAINER_TEST)

test-int-img: build-test-img
	docker run --rm -d --name=$(CONTAINER_TEST) abftest:latest
	docker exec $(CONTAINER_TEST) make test-int
	docker stop $(CONTAINER_TEST)

test-e2e-img: build-test-img
	docker run --rm -d --name=$(CONTAINER_TEST) abftest:latest
	docker exec $(CONTAINER_TEST) make test-e2e
	docker stop $(CONTAINER_TEST)

lint:
	golangci-lint run ./app/...

.PHONY: build build-img run run-img version test lint