# Makefile for a Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOGET=$(GOCMD) get
BINARY_DIR=bin
CMD_DIR=cmd
BINARY_NAME=$(BINARY_DIR)/barry

# Default target executed when no arguments are provided to make
all: build

# Build the project
build: $(BINARY_NAME)

$(BINARY_NAME): $(shell find $(CMD_DIR) -type f -name '*.go')
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_NAME) ./$(CMD_DIR)/...

# Format
format:
	$(GOFMT) ./...

# Run the application
run:
	$(BINARY_NAME)

# Clean up build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Run tests
test:
	$(GOTEST) -v ./...

.PHONY: all build clean deps test run
