BINARY_NAME := tuwi
PREFIX ?= /usr
BINDIR := $(PREFIX)/bin
CROSS_TARGETS ?= linux/amd64 linux/arm64 linux/386 linux/arm

.PHONY: install-deps build cross-build install install-local uninstall \
	run-dev run-debug run unit-test fuzz test lint format-check \
	format-write clean help

install-deps:
	@go mod download

build:
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .

cross-build:
	@failed=0; \
	for t in $(CROSS_TARGETS); do \
		os=$$(echo $$t | cut -d/ -f1); \
		arch=$$(echo $$t | cut -d/ -f2); \
		arm=$$(echo $$t | cut -s -d/ -f3); \
		printf '%-14s' "$$t"; \
		if out=$$(CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch GOARM=$$arm go build $(PKGS) 2>&1); then \
			echo OK; \
		else \
			echo FAIL; echo "$$out" | sed 's/^/    /'; failed=1; \
		fi; \
	done; \
	exit $$failed

install: build
	@sudo install -m 755 ./$(BINARY_NAME) /usr/bin
	@rm -f ./$(BINARY_NAME)
	@mkdir -p ~/.config/tuwi
	@touch ~/.config/tuwi/config.toml

install-local:
	@go install

uninstall:
	@sudo rm -f $(BINDIR)/$(BINARY_NAME)
	@rm -rf ~/.config/tuwi

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
	@echo "    make install        — Combile binary and move into $(BINDIR)" and create the config directory and blank config file in ~/.config/tuwi/config.toml
	@echo "    make uninstall      — Remove the binary from $(BINDIR)" and delete the ~/.config/tuwi directory
	@echo "    make run            — Build and execute the compiled binary"
	@echo "    make clean          — Remove clear cache, remove binary from $GOPATH/bin/ & $(BINDIR)"
	@echo ""
