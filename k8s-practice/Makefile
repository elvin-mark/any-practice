.PHONY: build test clean install

BIN_DIR := bin
BINARY := $(BIN_DIR)/cka

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BINARY) ./cmd/cka

test:
	go test -v ./...

install: build
	cp $(BINARY) /usr/local/bin/cka 2>/dev/null || cp $(BINARY) $(HOME)/.local/bin/cka

clean:
	rm -rf $(BIN_DIR)
	rm -f cka-exam-report-*.md
