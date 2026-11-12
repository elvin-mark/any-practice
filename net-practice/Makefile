.PHONY: all build test clean

BINARY_NAME := netp
BUILD_DIR := bin

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/netp

test:
	go test ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
