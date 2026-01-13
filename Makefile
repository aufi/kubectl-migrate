# Makefile for kubectl-migrate

# Version and Build Info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Binary Info
BINARY_NAME = kubectl-migrate
MAIN_PACKAGE = .
BUILD_DIR = bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Linker flags to inject version info
LDFLAGS = -ldflags "\
	-X 'github.com/konveyor/kubectl-migrate/internal/buildinfo.Version=$(VERSION)' \
	-X 'github.com/konveyor/kubectl-migrate/internal/buildinfo.Commit=$(COMMIT)' \
	-X 'github.com/konveyor/kubectl-migrate/internal/buildinfo.BuildDate=$(BUILD_DATE)'"

# Platforms for cross-compilation
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: all build build-all clean test deps install help

all: build

## help: Display this help message
help:
	@echo "kubectl-migrate Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## build: Build the binary for the current platform
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## build-all: Build binaries for all platforms
build-all: clean
	@echo "Building for all platforms..."
	@$(foreach platform,$(PLATFORMS), \
		$(eval OS = $(word 1,$(subst /, ,$(platform)))) \
		$(eval ARCH = $(word 2,$(subst /, ,$(platform)))) \
		$(eval OUTPUT = $(BUILD_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)) \
		$(if $(filter windows,$(OS)),$(eval OUTPUT := $(OUTPUT).exe)) \
		echo "Building for $(OS)/$(ARCH)..." && \
		GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) $(LDFLAGS) -o $(OUTPUT) $(MAIN_PACKAGE) && \
		echo "  ✓ $(OUTPUT)" || exit 1; \
	)
	@echo "All builds complete!"

## clean: Remove built binaries
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated!"

## install: Build and install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Installation complete!"
	@echo "You can now use: kubectl migrate <command>"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## lint: Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

## check: Run fmt, vet, and test
check: fmt vet test
	@echo "All checks passed!"
