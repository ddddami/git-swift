BINARY_NAME=git-swift
BUILD_DIR=bin

.DEFAULT_GOAL := build

.PHONY: build clean run rebuild

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/git-swift

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned."

run:
	@go run ./cmd/git-swift

rebuild: clean build
