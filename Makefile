all: lint test
PHONY: test coverage lint golint clean vendor local-dev-databases docker-up docker-down integration-test unit-test
GOOS=linux
OS_NAME := $(shell uname -s | tr A-Z a-z)
GOLANGCILINTCMD :=$(if ifeq darwin $(OS_NAME), docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.42 golangci-lint, golangci-lint)
GOCMD :=$(if ifeq darwin $(OS_NAME), docker run --rm -v $(shell pwd):/app -w /app golang:1.17 go, go)

test: | unit-test

unit-test: | lint
	@echo Running unit tests...
	@$(GOCMD) test -cover -short -tags testtools ./...

coverage:
	@echo Generating coverage report...
	@$(GOCMD) test ./... -race -coverprofile=coverage.out -covermode=atomic -tags testtools -p 1
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out

lint: golint

golint: | vendor
	@echo Linting Go files...
	@$(GOLANGCILINTCMD) run

clean: docker-clean
	@echo Cleaning...
	@rm -rf ./dist/
	@rm -rf coverage.out
	@go clean -testcache

vendor:
	@go mod download
	@go mod tidy -compat=1.17
