.PHONY: all build test clean

BINARY_NAME := algo
BUILD_DIR := bin

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/algo

test:
	go test ./internal/astcheck ./internal/exercise ./internal/progress ./internal/ui

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
