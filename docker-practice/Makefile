.PHONY: all build test clean install

BINARY_NAME := dockp
BUILD_DIR := bin

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/dockp

test:
	go test -v ./internal/...

clean:
	rm -rf $(BUILD_DIR)
	go clean

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME) 2>/dev/null || cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
