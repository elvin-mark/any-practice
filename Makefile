.PHONY: all build build-algo build-lp build-cka test clean help

BIN_DIR := bin

all: build

## build: Compiles all CLI/TUI tools into their subfolders and root bin/
build: build-algo build-lp build-cka
	@echo "✓ All Any-Practice tools successfully built!"

build-algo:
	@echo "Building algo-practice/bin/algo..."
	@cd algo-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp algo-practice/bin/algo $(BIN_DIR)/algo

build-lp:
	@echo "Building linux-practice/bin/lp..."
	@cd linux-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp linux-practice/bin/lp $(BIN_DIR)/lp

build-cka:
	@echo "Building k8s-practice/bin/cka..."
	@cd k8s-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp k8s-practice/bin/cka $(BIN_DIR)/cka

## test: Runs test suites across all three tools
test:
	@echo "Testing algo-practice..."
	@cd algo-practice && $(MAKE) test
	@echo "Testing linux-practice..."
	@cd linux-practice && go test ./internal/state ./internal/scenario ./internal/ui
	@echo "Testing k8s-practice..."
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
	@cd algo-practice && $(MAKE) clean
	@cd linux-practice && $(MAKE) clean
	@cd k8s-practice && $(MAKE) clean
	@rm -f *.test *.out coverage.txt
	@echo "✓ Clean complete!"

## help: Display this help message
help:
	@echo "Any-Practice - Unified Terminal Training Suite"
	@echo "Available commands:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //' | awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
