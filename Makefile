.PHONY: all build build-algo build-lp build-cka test clean help

BIN_DIR := bin

all: build

## build: Compiles all CLI/TUI tools (algo, lp, cka) into root bin/
build: build-algo build-lp build-cka
	@echo "✓ All Any-Practice tools successfully built into $(BIN_DIR)/"

build-algo:
	@mkdir -p $(BIN_DIR)
	@echo "Building algo..."
	@cd algo-practice && go build -o ../$(BIN_DIR)/algo ./cmd/algo
	@mkdir -p algo-practice/bin && cp $(BIN_DIR)/algo algo-practice/bin/algo

build-lp:
	@mkdir -p $(BIN_DIR)
	@echo "Building lp (linux-practice)..."
	@cd linux-practice && go build -o ../$(BIN_DIR)/lp ./cmd/lp
	@mkdir -p linux-practice/bin && cp $(BIN_DIR)/lp linux-practice/bin/lp

build-cka:
	@mkdir -p $(BIN_DIR)
	@echo "Building cka (k8s-practice)..."
	@cd k8s-practice && go build -o ../$(BIN_DIR)/cka ./cmd/cka
	@mkdir -p k8s-practice/bin && cp $(BIN_DIR)/cka k8s-practice/bin/cka

## test: Runs test suites across all three tools
test:
	@echo "Testing algo-practice (internal packages)..."
	@cd algo-practice && go test ./internal/astcheck ./internal/exercise ./internal/progress ./internal/ui
	@echo "Testing linux-practice (internal packages)..."
	@cd linux-practice && go test ./internal/state ./internal/scenario ./internal/ui
	@echo "Testing k8s-practice (internal packages)..."
	@cd k8s-practice && go test ./internal/exam ./internal/scenario ./internal/ui
	@echo "✓ All test suites passed!"

## vet: Runs go vet across all projects
vet:
	@echo "Running go vet..."
	@cd algo-practice && go vet ./...
	@cd linux-practice && go vet ./...
	@cd k8s-practice && go vet ./...
	@echo "✓ Go vet passed!"

## clean: Removes build binaries and transient test files
clean:
	@echo "Cleaning binaries..."
	@rm -rf $(BIN_DIR)
	@rm -f algo-practice/bin/algo algo-practice/algo
	@rm -f linux-practice/bin/lp
	@rm -f k8s-practice/bin/cka
	@rm -f *.test *.out coverage.txt
	@echo "✓ Clean complete!"

## help: Display this help message
help:
	@echo "Any-Practice - Unified Terminal Training Suite"
	@echo "Available commands:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //' | awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
