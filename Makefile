BINARY=packer-builder-sandwich
PACKAGES=$(shell glide novendor)
SOURCE_FILES=$(shell find . -name '*.go' -not -path '*vendor*')
VERSION=dev

.PHONY: all build check clean coverage fmt help lint test vet

all: check build test ## run fmt, vet, lint, build the binaries and run the tests

check: fmt vet lint ## run fmt, vet, lint

vet: ## run go vet
	@echo "Running $@"
	@test -z "$$(go vet ${PACKAGES} 2>&1 | tee /dev/stderr)"

fmt: ## run go fmt
	@echo "Running $@"
	@gofmt -s -l -w ${SOURCE_FILES}

build: ## build the go packages
	@echo "Running $@"
	@go build -i -ldflags "-X main.Version=${VERSION}" -o bin/${BINARY} .

build-linux: ## build the go packages for Linux (useful to copy the binary into docker)
	@echo "Running $@"
	@GOOS=linux GOARCH=amd64 go build -i -ldflags "-X main.Version=${VERSION}" -o bin/${BINARY} .

test: ## run test
	@echo "Running $@"
	@go test ${PACKAGES}

coverage: ## run tests with coverage metrics
	@echo "Running $@"
	@go test -cover ${PACKAGES}

clean: ## clean up binaries
	@echo "Running $@"
	@rm -rf bin

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort