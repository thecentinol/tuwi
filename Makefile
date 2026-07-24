BINARY_NAME := tuwi
PREFIX ?= /usr
BINDIR := $(PREFIX)/bin

.PHONY: install-deps build install install-local uninstall \
	run-dev run-debug run unit-test fuzz test lint format-check \
	format-write clean help

install-deps:
	@go mod download

build:
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .

install: build
	@sudo install -m 755 ./$(BINARY_NAME) /usr/bin
	@rm -f ./$(BINARY_NAME)

install-local:
	@go install

uninstall:
	@sudo rm -f $(BINDIR)/$(BINARY_NAME)

run-dev:
	@go run main.go

run-debug:
	@go run main.go -debug

run: build
	@./$(BINARY_NAME)

unit-test:
	@go test ./... -short

fuzz:
	@go test -fuzz=FuzzDetermineSecurityType -fuzztime=30s -v ./internal/networkmanager/...

test: unit-test

lint:
	@golangci-lint run

format-check:
	@gofumpt -l -extra .

format-write:
	@gofumpt -l -extra -w .

clean:
	@go clean
	@rm -rf ./dist
	@rm -f ~/go/bin/$(BINARY_NAME)
	@rm -f $(BINARY_NAME)

help:
	@echo ""
	@echo "  Setup:"
	@echo "    make install-deps   — Install project dependencies"
	@echo ""
	@echo "  Local development:"
	@echo "    make run-dev        — Run from source code"
	@echo "    make run-debug      — Run with debug"
	@echo "    make unit-test      — Run short unit tests"
	@echo "    make fuzz           — Run fuzz tests"
	@echo "    make test           — Run full test suite"
	@echo "    make lint           — Run linter"
	@echo "    make format-check   — Check the codebases formatting"
	@echo "    make format-write   — Fix formatting issues"
	@echo ""
	@echo "  Builds:"
	@echo "    make build          — Compile the binary into repo root"
	@echo "    make install        — Combile binary and move into $(BINDIR)"
	@echo "    make uninstall      — Remove the binary from $(BINDIR)"
	@echo "    make run            — Build and execute the compiled binary"
	@echo "    make clean          — Remove clear cache, remove binary from $GOPATH/bin/ & $(BINDIR)"
	@echo ""
