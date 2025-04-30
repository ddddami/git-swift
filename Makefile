BINARY_NAME=gitch
BUILD_DIR=bin

.DEFAULT_GOAL := build

.PHONY: build clean run rebuild

fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gitch

clean:
	@rm -rf $(BUILD_DIR)
	@echo "Cleaned."

run:
	@go run ./cmd/gitch

rebuild: clean build
