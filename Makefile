.PHONY: all build build-algo build-lp build-cka build-ctf build-dockp build-gitp build-sqlp test clean help

BIN_DIR := bin

all: build

## build: Compiles all CLI/TUI tools into their subfolders and root bin/
build: build-algo build-lp build-cka build-ctf build-dockp build-gitp build-sqlp
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

build-ctf:
	@echo "Building ctf-practice/bin/ctf..."
	@cd ctf-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp ctf-practice/bin/ctf $(BIN_DIR)/ctf

build-dockp:
	@echo "Building docker-practice/bin/dockp..."
	@cd docker-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp docker-practice/bin/dockp $(BIN_DIR)/dockp

build-gitp:
	@echo "Building git-practice/bin/gitp..."
	@cd git-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp git-practice/bin/gitp $(BIN_DIR)/gitp

build-sqlp:
	@echo "Building sql-practice/bin/sqlp..."
	@cd sql-practice && $(MAKE) build
	@mkdir -p $(BIN_DIR)
	@cp sql-practice/bin/sqlp $(BIN_DIR)/sqlp

## test: Runs test suites across all tools
test:
	@echo "Testing algo-practice..."
	@cd algo-practice && $(MAKE) test
	@echo "Testing linux-practice..."
	@cd linux-practice && go test ./internal/state ./internal/scenario ./internal/progress ./internal/ui
	@echo "Testing k8s-practice..."
	@cd k8s-practice && go test ./internal/exam ./internal/scenario ./internal/progress ./internal/ui
	@echo "Testing ctf-practice..."
	@cd ctf-practice && $(MAKE) test
	@echo "Testing docker-practice..."
	@cd docker-practice && $(MAKE) test
	@echo "Testing git-practice..."
	@cd git-practice && $(MAKE) test
	@echo "Testing sql-practice..."
	@cd sql-practice && $(MAKE) test
	@echo "✓ All test suites passed!"

## vet: Runs go vet across all projects
vet:
	@echo "Running go vet..."
	@cd algo-practice && go vet ./...
	@cd linux-practice && go vet ./...
	@cd k8s-practice && go vet ./...
	@cd ctf-practice && go vet ./...
	@cd docker-practice && go vet ./...
	@cd git-practice && go vet ./...
	@cd sql-practice && go vet ./...
	@echo "✓ Go vet passed!"

## clean: Removes build binaries and transient test files
clean:
	@echo "Cleaning binaries..."
	@rm -rf $(BIN_DIR)
	@cd algo-practice && $(MAKE) clean
	@cd linux-practice && $(MAKE) clean
	@cd k8s-practice && $(MAKE) clean
	@cd ctf-practice && $(MAKE) clean
	@cd docker-practice && $(MAKE) clean
	@cd git-practice && $(MAKE) clean
	@cd sql-practice && $(MAKE) clean
	@rm -f *.test *.out coverage.txt
	@echo "✓ Clean complete!"

## help: Display this help message
help:
	@echo "Any-Practice - Unified Terminal Training Suite"
	@echo "Available commands:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //' | awk 'BEGIN {FS = ": "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
