BINARY   = agentspec
PKG      = ./cmd/agentspec
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test install clean lint fmt help

## build: Compile the binary
build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) $(PKG)

## test: Run all tests
test:
	go test -race -v ./...

## install: Install the binary to $GOPATH/bin
install:
	go install $(PKG)

## clean: Remove build artifacts
clean:
	rm -f $(BINARY)
	rm -rf generated/

## lint: Run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## fmt: Format all Go files
fmt:
	gofmt -s -w .

## example-validate: Validate the example spec
example-validate: build
	./$(BINARY) validate examples/task-management/spec.yaml

## example-generate: Generate code from example spec
example-generate: build
	./$(BINARY) generate --target backend -f examples/task-management/spec.yaml -o /tmp/agentspec-gen/backend
	./$(BINARY) generate --target frontend -f examples/task-management/spec.yaml -o /tmp/agentspec-gen/frontend
	./$(BINARY) generate --target test -f examples/task-management/spec.yaml -o /tmp/agentspec-gen/test

## example-context: Generate agent context from example spec
example-context: build
	./$(BINARY) context --role backend -f examples/task-management/spec.yaml

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/  /'
