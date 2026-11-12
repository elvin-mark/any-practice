.PHONY: all build test clean install

BINARY_NAME=lp
BUILD_DIR=bin
MODULE=github.com/elvin-mark/linux-practice

all: test build

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/lp

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
