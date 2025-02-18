.PHONY: all
all: tidy build unit-test lint 

.PHONY: build
build:
	echo "Build Go SDK"
	go build -v ./...

.PHONY: tidy
tidy:
	echo "Tidy dependency modules"
	go mod tidy

.PHONY: download
download:
	echo "Download dependency modules"
	go mod download

.PHONY: unit-test
unit-test: 
	echo "Running unit tests for service"
	go clean -testcache
	go test ./... -vet=off $(ARGS)

.PHONY: test-coverage-report
test-coverage-report: 
	go test ./... -cover -coverprofile coverage.out -covermode count
	go tool cover -html=coverage.out

lint-install: 
	echo "Installing golangci-lint"
	brew install golangci-lint
	brew upgrade golangci-lint	

.PHONY: lint
lint:
	echo "Running checks for service"
	golangci-lint run ./...