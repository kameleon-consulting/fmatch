# fmatch — Makefile
# Targets: build, test, lint, cross-compile
# Requires: go 1.24+, golangci-lint (for lint target)

BINARY_NAME := fmatch
MODULE      := github.com/mlabate/fmatch
BUILD_DIR   := dist
GO          := go
GOFLAGS     :=

# Build info injected at link time
VERSION     := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS     := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION) -s -w"

# Cross-compile targets
PLATFORMS   := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build test lint clean cross-compile help

## all: build the binary (default)
all: build

## build: compile the binary for the current OS/ARCH
build:
	@echo "==> Building $(BINARY_NAME) ($(VERSION))..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "    Done: ./$(BINARY_NAME)"

## test: run all unit tests with race detector
test:
	@echo "==> Running tests..."
	$(GO) test -race -count=1 ./...

## lint: run static analysis (requires golangci-lint)
lint:
	@echo "==> Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint not found. Install it: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	}
	golangci-lint run ./...

## cross-compile: build binaries for all supported platforms
cross-compile: $(PLATFORMS)

$(PLATFORMS):
	$(eval OS   := $(word 1, $(subst /, ,$@)))
	$(eval ARCH := $(word 2, $(subst /, ,$@)))
	$(eval EXT  := $(if $(filter windows,$(OS)),.exe,))
	@echo "==> Building $(BINARY_NAME)-$(OS)-$(ARCH)$(EXT) ($(VERSION))..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GOFLAGS) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(EXT) .

## clean: remove build artifacts
clean:
	@echo "==> Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "    Done."

## help: show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
