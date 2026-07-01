BINARY_NAME := tuwi

.PHONY: all build clean help install run run-debug run-dev test

all: build

build:
	@echo "Building $(BINARY_NAME) binary..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .

install:
	@go mod download

run-dev:
	@go run main.go

run-debug:
	@go run main.go -debug

run: build
	@./$(BINARY_NAME)

unit-test:
	@go test ./... -short

fuzz:
	go test -fuzz=FuzzDetermineSecurityType -fuzztime=30s -v ./internal/networkmanager/...

test: unit-test

clean:
	@go clean
	@rm -rf $(BINARY_NAME)

help:
	@echo ""
	@echo "  Setup:"
	@echo "    make install        - Install project dependencies listed in go.mod"
	@echo ""
	@echo "  Local Development:"
	@echo "    make run-dev        - Run from source code"
	@echo "    make run-debug      - Run with debug"
	@echo "    make test           - Run full test suite"
	@echo "    make unit-test      - Run short unit tests"
	@echo "    make fuzz           - Run fuzz tests"
	@echo ""
	@echo "  Local Builds:"
	@echo "    make build          - Compile the binary"
	@echo "    make run            - Build and execute the compiled binary"
	@echo "    make clean          - Remove binary and clear cache"
	@echo ""
