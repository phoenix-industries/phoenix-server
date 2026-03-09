SRC := $(wildcard cmd/phoenix-server/*.go pkg/**/*.go internal/**/*.go assets/*)

run: build
	./bin/phoenix-server
.PHONY: run

## build: build the application
build: /bin/phoenix-server
.PHONY: build

bin/phoenix-server: $(SRC) scripts/build
	./scripts/build phoenix-server

## run/live: run the application with reloading on file changes
run/live: .air.toml
	go run github.com/air-verse/air@latest
.PHONY: run/live

## audit: run quality control checks
audit: test
	go mod tidy -diff
	go mod verify
	test -z "$(shell gofmt -l .)"
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@command -v bodyclose >/dev/null 2>&1 || go install github.com/timakin/bodyclose@latest
	go vet -vettool="$$(which bodyclose)" ./...
	go run github.com/client9/misspell/cmd/misspell@latest -locale="US" -error -source="text" **/*
.PHONY: audit

## test: run all tests
test:
	go test -v -race -buildvcs ./...
.PHONY: test

## test/cover: run all tests and display coverage
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out
.PHONY: test/cover

## upgradeable: list direct dependencies that have upgrades available
upgradeable:
	@go run github.com/oligot/go-mod-upgrade@latest
.PHONY: upgradeable

## tidy: tidy modfiles and format .go files
tidy: fmt deps
.PHONY: tidy

## fmt: format .go files
fmt:
	go fmt ./...
	find . -name \*.go -not -path .git -not -path bin -exec go run golang.org/x/tools/cmd/goimports@latest -w {} \;
.PHONY: fmt

## deps: install and verify dependencies
deps:
	go get -u ./...
	go mod verify
	go mod tidy -v
.PHONY: deps

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: help
